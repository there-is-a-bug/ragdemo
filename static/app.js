const API_BASE = "";
const UPLOAD_ENDPOINT = "/api/upload";
const RETRIEVE_ENDPOINT = "/api/retrieve";

function generateUID() {
    function generateUID() {
        return Math.floor(
            100000 + Math.random() * 900000
        );
    }
    const UID = generateUID();
    return Number(UID);
}
const UID = generateUID();
document.getElementById("uid").textContent = UID;
const contentInput =
    document.getElementById("content");
const chunkInput =
    document.getElementById("chunk");
const knnScoreInput =
    document.getElementById("knnScore");
const topKInput =
    document.getElementById("topK");
const uploadBtn =
    document.getElementById("uploadBtn");
const uploadStatus =
    document.getElementById("uploadStatus");
const chatForm =
    document.getElementById("chatForm");
const queryInput =
    document.getElementById("query");
const messages =
    document.getElementById("messages");
uploadBtn.addEventListener(
    "click",
    async () => {
        const content =
            contentInput.value;
        if (!content.trim()) {
            setStatus(
                "请输入文档内容。",
                "error"
            );
            return;
        }
        const chunk =
            Number(chunkInput.value);
        const body = {
            user_id: UID,
            content: content,
            chunk: chunk
        };
        uploadBtn.disabled = true;
        setStatus(
            "上传中...",
            ""
        );
        try {
            const response =
                await fetch(
                    API_BASE + UPLOAD_ENDPOINT,
                    {
                        method: "POST",
                        headers: {
                            "Content-Type":
                                "application/json"
                        },
                        body:
                            JSON.stringify(body)
                    }
                );
            const data =
                await parseResponse(response);
            if (!response.ok) {
                throw new Error(
                    extractError(
                        data,
                        response.status
                    )
                );
            }
            setStatus(
                "文档上传成功。",
                "success"
            );
        } catch (error) {
            setStatus(
                "上传失败：" +
                error.message,
                "error"
            );
        } finally {
            uploadBtn.disabled = false;
        }
    }
);

chatForm.addEventListener(
    "submit",
    async (event) => {
        event.preventDefault();
        const query =
            queryInput.value.trim();
        if (!query) {
            return;
        }
        removeEmptyState();
        appendUserMessage(query);
        queryInput.value = "";
        const loading =
            appendLoading();
        const knnScore =
            Number(knnScoreInput.value);
        const topK =
            Number(topKInput.value);
        if (
            !Number.isInteger(knnScore) ||
            knnScore < 1 ||
            knnScore > 100
        ) {
            loading.remove();
            appendError(
                "KNN Score 必须是 1-100 的整数。"
            );
            return;
        }
        if (
            !Number.isInteger(topK) ||
            topK < 1
        ) {
            loading.remove();
            appendError(
                "TopK 必须是大于等于 1 的整数。"
            );
            return;
        }
        const body = {
            user_id: UID,
            content: query,
            knn_score: knnScore,
            top_k: topK
        };
        try {
            const response =
                await fetch(
                    API_BASE + RETRIEVE_ENDPOINT,
                    {
                        method: "POST",
                        headers: {
                            "Content-Type":
                                "application/json"
                        },
                        body:
                            JSON.stringify(body)
                    }
                );
            const data =
                await parseResponse(response);
            loading.remove();
            if (!response.ok) {
                throw new Error(
                    extractError(
                        data,
                        response.status
                    )
                );
            }
            renderResults(
                data,
                {
                    query,
                    knnScore,
                    topK
                }
            );
        } catch (error) {
            loading.remove();
            appendError(
                "召回失败：" +
                error.message
            );
        }
    }
);

async function parseResponse(response) {
    const text =
        await response.text();
    if (!text) {
        return null;
    }
    try {
        return JSON.parse(text);
    } catch {
        return text;
    }
}

function extractError(
    data,
    status
) {
    if (
        typeof data === "string"
    ) {
        return (
            data ||
            `HTTP ${status}`
        );
    }
    if (
        data &&
        typeof data === "object"
    ) {
        return (
            data.message ||
            data.error ||
            `HTTP ${status}`
        );
    }
    return `HTTP ${status}`;
}

function removeEmptyState() {
    const empty =
        messages.querySelector(".empty");
    if (empty) {
        empty.remove();
    }
}

function appendUserMessage(text) {
    const wrapper =
        document.createElement("div");
    wrapper.className =
        "message user-message";
    const bubble =
        document.createElement("div");
    bubble.className =
        "user-bubble";
    bubble.textContent =
        text;
    wrapper.appendChild(bubble);
    messages.appendChild(wrapper);
    scrollToBottom();
}

function appendLoading() {
    const wrapper =
        document.createElement("div");
    wrapper.className =
        "message";
    const box =
        document.createElement("div");
    box.className =
        "result-message loading";
    box.textContent =
        "正在召回...";
    wrapper.appendChild(box);
    messages.appendChild(wrapper);
    scrollToBottom();
    return wrapper;
}

function appendError(message) {
    const wrapper =
        document.createElement("div");
    wrapper.className =
        "message";
    const box =
        document.createElement("div");
    box.className =
        "error-box";
    box.textContent =
        message;
    wrapper.appendChild(box);
    messages.appendChild(wrapper);
    scrollToBottom();
}

function renderResults(
    data,
    meta
) {
    const results =
        normalizeResults(data);
    const wrapper =
        document.createElement("div");
    wrapper.className =
        "message";
    const box =
        document.createElement("div");
    box.className =
        "result-message";
    const head =
        document.createElement("div");
    head.className =
        "result-head";
    const title =
        document.createElement("div");
    title.className =
        "result-title";
    title.textContent =
        `召回结果（${results.length} 条）`;
    const resultMeta =
        document.createElement("div");
    resultMeta.className =
        "result-meta";
    resultMeta.textContent =
        `KNN ≥ ${meta.knnScore} · TopK ${meta.topK}`;
    head.appendChild(title);
    head.appendChild(resultMeta);
    box.appendChild(head);
    if (results.length === 0) {
        const empty =
            document.createElement("div");
        empty.className =
            "loading";
        empty.textContent =
            "没有召回到结果。";
        box.appendChild(empty);
    }
    else {
        results.forEach(
            (item, index) => {
                box.appendChild(
                    createResultCard(
                        item,
                        index
                    )
                );
            }
        );
    }
    wrapper.appendChild(box);
    messages.appendChild(wrapper);
    scrollToBottom();
}

function createResultCard(
    item,
    index
) {
    const card =
        document.createElement("div");
    card.className =
        "result-card";
    const rank =
        document.createElement("div");
    rank.className =
        "result-row";
    const rankLabel =
        document.createElement("span");
    rankLabel.className =
        "result-label";
    rankLabel.textContent =
        "排名";
    const rankValue =
        document.createElement("strong");
    rankValue.textContent =
        `#${index + 1}`;
    rank.appendChild(rankLabel);
    rank.appendChild(rankValue);
    const score =
        document.createElement("div");
    score.className =
        "result-row";
    const scoreLabel =
        document.createElement("span");
    scoreLabel.className =
        "result-label";
    scoreLabel.textContent =
        "评分";
    const scoreValue =
        document.createElement("strong");
    scoreValue.className =
        "score";
    scoreValue.textContent =
        formatScore(item.score);
    score.appendChild(scoreLabel);
    score.appendChild(scoreValue);
    const source =
        document.createElement("div");
    source.className =
        "result-row";
    const sourceLabel =
        document.createElement("span");
    sourceLabel.className =
        "result-label";
    sourceLabel.textContent =
        "召回源";
    const sourceValue =
        document.createElement("span");
    sourceValue.className =
        "source";
    sourceValue.textContent =
        item.source ?? "";
    source.appendChild(sourceLabel);
    source.appendChild(sourceValue);
    const content =
        document.createElement("div");
    content.className =
        "result-content";
    content.textContent =
        item.content ?? "";
    card.appendChild(rank);
    card.appendChild(score);
    card.appendChild(source);
    card.appendChild(content);
    return card;
}

function normalizeResults(data) {
    if (Array.isArray(data)) {

        return data;

    }
    if (
        data &&
        Array.isArray(data.data)
    ) {
        return data.data;
    }
    if (
        data &&
        Array.isArray(data.results)
    ) {
        return data.results;
    }
    if (
        data &&
        data.content !== undefined
    ) {
        return [data];
    }
    return [];
}

function formatScore(score) {
    const value =
        Number(score);
    if (!Number.isFinite(value)) {
        return "-";
    }
    return value.toFixed(4);
}

function setStatus(
    text,
    type
) {
    uploadStatus.textContent =
        text;
    uploadStatus.className =
        "status";
    if (type) {
        uploadStatus.classList.add(
            type
        );
    }
}

function scrollToBottom() {
    messages.scrollTop =
        messages.scrollHeight;
}