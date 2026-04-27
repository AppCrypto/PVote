const appState = {
  activeTab: "initiator",
  snapshot: null,
};

document.addEventListener("DOMContentLoaded", async () => {
  bindTabs();
  bindForms();
  bindTallierActions();
  await refreshState();
});

function bindTabs() {
  document.querySelectorAll(".role-tab").forEach((button) => {
    button.addEventListener("click", () => {
      appState.activeTab = button.dataset.tab;
      renderTabs();
    });
  });
}

function bindForms() {
  document.getElementById("setup-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const payload = {
      numTalliers: Number(document.getElementById("setup-talliers").value),
      numCandidates: Number(document.getElementById("setup-candidates").value),
      threshold: Number(document.getElementById("setup-threshold").value),
      rangeMin: Number(document.getElementById("setup-range-min").value),
      rangeMax: Number(document.getElementById("setup-range-max").value),
    };
    await postJSON("/api/setup", payload);
    appState.activeTab = "initiator";
    renderTabs();
  });

  document.getElementById("vote-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const alias = document.getElementById("vote-alias").value.trim();
    const scores = Array.from(document.querySelectorAll("[data-score-input]")).map((input) => Number(input.value));
    await postJSON("/api/voters", { alias, scores });
    document.getElementById("vote-alias").value = "";
  });
}

function bindTallierActions() {
  document.getElementById("tallier-grid").addEventListener("click", async (event) => {
    const button = event.target.closest("[data-decrypt-id]");
    if (!button) {
      return;
    }
    await postJSON(`/api/talliers/${button.dataset.decryptId}/decrypt`, {});
  });

  document.getElementById("finalize-button").addEventListener("click", async () => {
    await postJSON("/api/tally/finalize", {});
  });
}

async function refreshState() {
  const response = await fetch("/api/state");
  const payload = await response.json();
  appState.snapshot = payload.state;
  render();
}

async function postJSON(url, payload) {
  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const data = await response.json();
  if (data.state) {
    appState.snapshot = data.state;
    render();
  }

  if (!response.ok) {
    flash(data.error || "Request failed", true);
    return;
  }

  flash(data.message || "Updated");
}

function render() {
  if (!appState.snapshot) {
    return;
  }

  renderTabs();
  renderOverview();
  renderSetupForm();
  renderParams();
  renderVoteComposer();
  renderVoters();
  renderTalliers();
  renderTally();
  renderTimeline();
}

function renderTabs() {
  document.querySelectorAll(".role-tab").forEach((button) => {
    button.classList.toggle("is-active", button.dataset.tab === appState.activeTab);
  });

  document.querySelectorAll(".tab-panel").forEach((panel) => {
    panel.classList.toggle("is-active", panel.id === `panel-${appState.activeTab}`);
  });
}

function renderOverview() {
  const { meta, overview } = appState.snapshot;
  document.getElementById("phase-chip").textContent = phaseTitle(meta.phase);
  document.getElementById("session-chip").textContent = meta.sessionId;
  document.getElementById("range-chip").textContent = `Range ${overview.rangeLabel}`;

  document.getElementById("metric-votes").textContent = String(overview.voteCount);
  document.getElementById("metric-verified").textContent = `${overview.verifiedVotes}`;
  document.getElementById("metric-talliers").textContent = `${overview.decryptionCount} / ${overview.threshold}`;
  document.getElementById("metric-tally").textContent = appState.snapshot.tally ? "Recovered" : phaseTitle(meta.phase);
}

function renderSetupForm() {
  const config = appState.snapshot.initiator.config;
  document.getElementById("setup-talliers").value = config.numTalliers;
  document.getElementById("setup-candidates").value = config.numCandidates;
  document.getElementById("setup-threshold").value = config.threshold;
  document.getElementById("setup-range-min").value = config.rangeMin;
  document.getElementById("setup-range-max").value = config.rangeMax;
}

function renderParams() {
  const params = appState.snapshot.initiator.publicParams;
  const paramGrid = document.getElementById("param-grid");
  paramGrid.innerHTML = [
    paramCard("G0", params.g0),
    paramCard("H0", params.h0),
    paramCard("G1", params.g1),
    paramCard("PKI", params.pki),
  ].join("");

  document.getElementById("tallier-keys").innerHTML = appState.snapshot.initiator.tallierKeys
    .map((item) => `<span class="token">T${item.id} · ${item.publicKey}</span>`)
    .join("");

  document.getElementById("aggregate-c").innerHTML = appState.snapshot.initiator.aggregate.encryptedShares
    .map((item, index) => `<span class="token token-soft">C${index + 1} · ${item}</span>`)
    .join("");

  document.getElementById("aggregate-u").innerHTML = appState.snapshot.initiator.aggregate.ballotCipher
    .map((item, index) => `<span class="token token-soft">U${index + 1} · ${item}</span>`)
    .join("");
}

function renderVoteComposer() {
  const { overview, initiator } = appState.snapshot;
  const scoreGrid = document.getElementById("score-grid");
  const rangeMin = initiator.config.rangeMin;
  const rangeMax = initiator.config.rangeMax;

  scoreGrid.innerHTML = initiator.candidateLabels
    .map((label, index) => {
      const midpoint = Math.floor((rangeMin + rangeMax) / 2);
      return `
        <label class="score-card">
          <span class="score-title">${label}</span>
          <input
            type="range"
            min="${rangeMin}"
            max="${rangeMax}"
            value="${midpoint}"
            data-score-input
            data-score-index="${index}"
          />
          <span class="score-value" data-score-value>${midpoint}</span>
        </label>
      `;
    })
    .join("");

  scoreGrid.querySelectorAll("[data-score-input]").forEach((input) => {
    const valueEl = input.parentElement.querySelector("[data-score-value]");
    input.addEventListener("input", () => {
      valueEl.textContent = input.value;
    });
  });

  document.getElementById("vote-submit").disabled = !overview.canVote;
}

function renderVoters() {
  const { voters } = appState.snapshot;
  const list = document.getElementById("voter-list");

  if (!voters.length) {
    list.innerHTML = emptyState("No voter has submitted a ballot yet.");
    return;
  }

  list.innerHTML = voters
    .slice()
    .reverse()
    .map(
      (voter) => `
        <article class="record-card">
          <div class="record-head">
            <div>
              <h3>${escapeHTML(voter.alias)}</h3>
              <p class="record-meta">#${voter.id} · ${voter.submittedAt}</p>
            </div>
            <div class="badge-row">
              <span class="badge ${voter.pvssVerified ? "badge-ok" : "badge-warn"}">PVSS ${voter.pvssVerified ? "OK" : "FAIL"}</span>
              <span class="badge ${voter.rangeVerified ? "badge-ok" : "badge-warn"}">ZKRP ${voter.rangeVerified ? "OK" : "FAIL"}</span>
            </div>
          </div>
          <p class="score-line">${formatScores(voter.scores, appState.snapshot.initiator.candidateLabels)}</p>
          <div class="token-list">
            ${voter.bindCommitments.map((item, index) => `<span class="token">Bind ${index + 1} · ${item}</span>`).join("")}
          </div>
          <details class="details">
            <summary>Show encrypted shares and ballot ciphertext</summary>
            <div class="details-grid">
              <div>${voter.encryptedShares.map((item, index) => `<span class="token token-soft">C${index + 1} · ${item}</span>`).join("")}</div>
              <div>${voter.ballotCipher.map((item, index) => `<span class="token token-soft">U${index + 1} · ${item}</span>`).join("")}</div>
            </div>
          </details>
        </article>
      `
    )
    .join("");
}

function renderTalliers() {
  const { talliers, overview } = appState.snapshot;
  const grid = document.getElementById("tallier-grid");

  grid.innerHTML = talliers
    .map(
      (tallier) => `
        <article class="tallier-card">
          <div class="record-head">
            <div>
              <h3>Tallier ${String(tallier.id).padStart(2, "0")}</h3>
              <p class="record-meta">${tallier.publicKey}</p>
            </div>
            <span class="badge ${tallier.verified ? "badge-ok" : "badge-neutral"}">
              ${tallier.hasDecrypted ? (tallier.verified ? "Verified" : "Published") : "Idle"}
            </span>
          </div>
          <p class="tallier-copy">
            ${tallier.hasDecrypted ? `Share ${tallier.share}<br />Proof ${tallier.proof}` : "Waiting to decrypt the aggregated ciphertext assigned to this tallier."}
          </p>
          <div class="tallier-actions">
            <button
              class="button-secondary"
              data-decrypt-id="${tallier.id}"
              ${!overview.canDecrypt || tallier.hasDecrypted ? "disabled" : ""}
            >
              ${tallier.hasDecrypted ? "Decrypted" : "Decrypt Share"}
            </button>
            <span class="record-meta">${tallier.decryptedAt || ""}</span>
          </div>
        </article>
      `
    )
    .join("");

  document.getElementById("finalize-button").disabled = !overview.canFinalize;
  document.getElementById("finalize-note").textContent = overview.canFinalize
    ? "Threshold satisfied. You can reconstruct the final tally now."
    : `Need ${overview.threshold} verified talliers before final reconstruction.`;
}

function renderTally() {
  const tally = appState.snapshot.tally;
  const board = document.getElementById("tally-board");

  if (!tally) {
    board.innerHTML = emptyState("The final tally will appear here once the decryption threshold is reached.");
    return;
  }

  const maxValue = Math.max(...tally.results.map((item) => item.total), 1);
  board.innerHTML = `
    <div class="badge-row">
      <span class="badge ${tally.verified ? "badge-ok" : "badge-warn"}">${tally.verified ? "Cross-check passed" : "Cross-check mismatch"}</span>
      <span class="record-meta">Finalized at ${tally.finalizedAt}</span>
    </div>
    <div class="result-list">
      ${tally.results
        .map(
          (item) => `
            <article class="result-card">
              <div class="result-head">
                <h3>${item.label}</h3>
                <p>${item.total}</p>
              </div>
              <div class="result-bar">
                <span style="width:${(item.total / maxValue) * 100}%"></span>
              </div>
              <p class="record-meta">Recovered commitment ${item.commitment}</p>
              <p class="record-meta">Plaintext cross-check ${item.expectedTotal}</p>
            </article>
          `
        )
        .join("")}
    </div>
  `;
}

function renderTimeline() {
  const timeline = document.getElementById("timeline");
  const events = appState.snapshot.timeline;

  if (!events.length) {
    timeline.innerHTML = emptyState("No events yet.");
    return;
  }

  timeline.innerHTML = events
    .map(
      (event) => `
        <article class="timeline-item">
          <div class="timeline-dot timeline-${event.role}"></div>
          <div class="timeline-copy">
            <div class="record-head">
              <h3>${escapeHTML(event.title)}</h3>
              <span class="record-meta">${event.at}</span>
            </div>
            <p class="record-meta timeline-role">${event.role}</p>
            <p>${escapeHTML(event.detail)}</p>
          </div>
        </article>
      `
    )
    .join("");
}

function flash(message, isError = false) {
  const toast = document.getElementById("toast");
  toast.hidden = false;
  toast.textContent = message;
  toast.classList.toggle("is-error", isError);
  toast.classList.add("is-visible");

  window.clearTimeout(flash.timer);
  flash.timer = window.setTimeout(() => {
    toast.classList.remove("is-visible");
    toast.hidden = true;
  }, 2800);
}

function phaseTitle(phase) {
  const map = {
    setup: "Setup",
    voting: "Voting",
    tallying: "Tallying",
    completed: "Completed",
  };
  return map[phase] || phase;
}

function paramCard(label, value) {
  return `
    <article class="stat-card">
      <p class="metric-label">${label}</p>
      <p class="token-line">${value}</p>
    </article>
  `;
}

function formatScores(scores, labels) {
  return scores
    .map((score, index) => `${labels[index]}: ${score}`)
    .join(" · ");
}

function emptyState(message) {
  return `<div class="empty-state">${message}</div>`;
}

function escapeHTML(value) {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}
