const appState = {
  activeTab: "initiator",
  snapshot: null,
};

document.addEventListener("DOMContentLoaded", async () => {
  bindTabs();
  bindForms();
  bindActions();
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
      initiatorEscrowEth: document.getElementById("setup-initiator-escrow").value,
      voterStakeEth: document.getElementById("setup-voter-stake").value,
      tallierStakeEth: document.getElementById("setup-tallier-stake").value,
      initiatorRewardPercent: Number(document.getElementById("setup-initiator-reward").value),
      voterRewardPercent: Number(document.getElementById("setup-voter-reward").value),
      tallierRewardPercent: Number(document.getElementById("setup-tallier-reward").value),
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

function bindActions() {
  document.getElementById("fund-escrow-button").addEventListener("click", async () => {
    await postJSON("/api/initiator/escrow", {});
  });

  document.getElementById("withdraw-initiator-button").addEventListener("click", async () => {
    await postJSON("/api/initiator/withdraw", {});
  });

  document.getElementById("voter-list").addEventListener("click", async (event) => {
    const withdraw = event.target.closest("[data-withdraw-voter-id]");
    if (withdraw) {
      await postJSON(`/api/voters/${withdraw.dataset.withdrawVoterId}/withdraw`, {});
    }
  });

  document.getElementById("tallier-grid").addEventListener("click", async (event) => {
    const stake = event.target.closest("[data-stake-id]");
    if (stake) {
      await postJSON(`/api/talliers/${stake.dataset.stakeId}/stake`, {});
      return;
    }

    const decrypt = event.target.closest("[data-decrypt-id]");
    if (decrypt) {
      await postJSON(`/api/talliers/${decrypt.dataset.decryptId}/decrypt`, {});
      return;
    }

    const withdraw = event.target.closest("[data-withdraw-tallier-id]");
    if (withdraw) {
      await postJSON(`/api/talliers/${withdraw.dataset.withdrawTallierId}/withdraw`, {});
    }
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
  const { meta, overview, chain, tally } = appState.snapshot;
  document.getElementById("phase-chip").textContent = phaseTitle(meta.phase);
  document.getElementById("session-chip").textContent = meta.sessionId;
  document.getElementById("range-chip").textContent = `Range ${overview.rangeLabel}`;

  document.getElementById("metric-votes").textContent = String(overview.voteCount);
  document.getElementById("metric-verified").textContent = `${overview.verifiedVotes}`;
  document.getElementById("metric-talliers").textContent = `${overview.decryptionCount} / ${overview.threshold}`;
  document.getElementById("metric-tally").textContent = chain.settled ? "Settled" : tally ? "Recovered" : phaseTitle(meta.phase);
}

function renderSetupForm() {
  const config = appState.snapshot.initiator.config;
  document.getElementById("setup-talliers").value = config.numTalliers;
  document.getElementById("setup-candidates").value = config.numCandidates;
  document.getElementById("setup-threshold").value = config.threshold;
  document.getElementById("setup-range-min").value = config.rangeMin;
  document.getElementById("setup-range-max").value = config.rangeMax;
  document.getElementById("setup-initiator-escrow").value = config.initiatorEscrowEth;
  document.getElementById("setup-voter-stake").value = config.voterStakeEth;
  document.getElementById("setup-tallier-stake").value = config.tallierStakeEth;
  document.getElementById("setup-initiator-reward").value = config.initiatorRewardPercent;
  document.getElementById("setup-voter-reward").value = config.voterRewardPercent;
  document.getElementById("setup-tallier-reward").value = config.tallierRewardPercent;
}

function renderParams() {
  const { initiator, chain, overview } = appState.snapshot;
  const params = initiator.publicParams;
  const paramGrid = document.getElementById("param-grid");
  paramGrid.innerHTML = [
    paramCard("G0", params.g0),
    paramCard("H0", params.h0),
    paramCard("G1", params.g1),
    paramCard("PKI", params.pki),
  ].join("");

  document.getElementById("tallier-keys").innerHTML = initiator.tallierKeys
    .map((item) => `<span class="token">T${item.id} · ${item.publicKey}</span>`)
    .join("");

  document.getElementById("aggregate-c").innerHTML = initiator.aggregate.encryptedShares
    .map((item, index) => `<span class="token token-soft">C${index + 1} · ${item}</span>`)
    .join("");

  document.getElementById("aggregate-u").innerHTML = initiator.aggregate.ballotCipher
    .map((item, index) => `<span class="token token-soft">U${index + 1} · ${item}</span>`)
    .join("");

  const chainStatus = document.getElementById("chain-status");
  chainStatus.textContent = chain.status;
  chainStatus.classList.toggle("is-error", !chain.available);

  const fundButton = document.getElementById("fund-escrow-button");
  fundButton.disabled = !overview.canFundInitiator;
  fundButton.textContent = chain.escrowFunded ? "Escrow Funded" : `Fund ${formatETH(chain.initiatorEscrowEth)} Escrow`;

  const withdrawInitiatorButton = document.getElementById("withdraw-initiator-button");
  withdrawInitiatorButton.disabled = !chain.initiator.canWithdraw;
  withdrawInitiatorButton.textContent = chain.initiator.withdrawn
    ? "Initiator Withdrawn"
    : chain.initiator.canWithdraw
      ? `Withdraw ${formatETH(chain.initiator.claimableEth)}`
      : "Withdraw Initiator Share";

  const escrowBadge = chain.available
    ? chain.escrowFunded
      ? statusBadge("Funded", "ok")
      : statusBadge("Needs funding", "warn")
    : statusBadge("Unavailable", "warn");
  const rewardBadge = chain.available
    ? chain.settled
      ? statusBadge("Settled", "ok")
      : statusBadge("Pending", "neutral")
    : statusBadge("No settlement", "warn");
  const voterStakeDetail = chain.available
    ? `${overview.remainingVoterSlots} Ganache voter wallets left`
    : "No Ganache voter wallets are active; ballots can run off-chain.";
  const tallierStakeDetail = chain.available
    ? "Each tallier must stake before decrypting."
    : "Talliers can decrypt off-chain while Ganache is unavailable.";
  const tallierStakeBadge = chain.available
    ? chain.settled
      ? statusBadge("Settlement closed", "neutral")
      : statusBadge("Open", "ok")
    : statusBadge("Off-chain only", "warn");

  document.getElementById("chain-overview").innerHTML = [
    ledgerCard("Ganache", chain.available ? "Connected" : "Unavailable", chain.rpcUrl, chain.available ? statusBadge("Live", "ok") : statusBadge("Offline", "warn")),
    ledgerCard("Contract", chain.contractAddress ? shortAddress(chain.contractAddress) : "Not deployed", chain.contractAddress || "Start Ganache with the README mnemonic, then rebuild the session."),
    ledgerCard("Escrow Balance", formatETH(chain.contractBalanceEth), `${formatETH(chain.totalEscrowEth)} total deposited`, escrowBadge),
    ledgerCard("Reward Pool", formatETH(chain.rewardPoolEth), chain.rewardSplit, rewardBadge),
    financeCard("Initiator", chain.initiator, chain.initiatorEscrowEth, "Reward escrow"),
    ledgerCard("Voter Stake", formatETH(chain.voterStakeEth), voterStakeDetail, chain.available ? statusBadge("Auto-staked", "neutral") : statusBadge("Off-chain only", "warn")),
    ledgerCard("Tallier Stake", formatETH(chain.tallierStakeEth), tallierStakeDetail, tallierStakeBadge),
  ].join("");
}

function renderVoteComposer() {
  const { overview, initiator, chain } = appState.snapshot;
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

  const voteSubmit = document.getElementById("vote-submit");
  voteSubmit.disabled = !overview.canVote;
  voteSubmit.textContent = !chain.available
    ? "Submit Off-chain Ballot"
    : chain.escrowFunded
      ? "Stake & Submit Ballot"
      : "Fund Escrow First";
  document.getElementById("vote-stake-note").textContent = chain.available
    ? `Each accepted ballot escrows ${formatETH(chain.voterStakeEth)} from a dedicated Ganache voter wallet. Remaining voter wallets: ${overview.remainingVoterSlots}.`
    : `Ganache escrow is unavailable. The vote can still run off-chain, but no on-chain stake or reward settlement is active.`;
}

function renderVoters() {
  const { voters, chain } = appState.snapshot;
  const list = document.getElementById("voter-list");

  if (!voters.length) {
    list.innerHTML = emptyState("No voter has submitted a ballot yet.");
    return;
  }

  list.innerHTML = voters
    .slice()
    .reverse()
    .map((voter) => {
      const withdrawLabel = voter.stake.withdrawn
        ? "Withdrawn"
        : voter.stake.canWithdraw
          ? `Withdraw ${formatETH(voter.stake.claimableEth)}`
          : "Withdraw After Settlement";
      const accountLabel = voter.stake.address
        ? `Ganache account ${shortAddress(voter.stake.address)}`
        : chain.available
          ? "Waiting for a Ganache voter account."
          : "No on-chain voter stake is active.";

      return `
        <article class="record-card">
          <div class="record-head">
            <div>
              <h3>${escapeHTML(voter.alias)}</h3>
              <p class="record-meta">#${voter.id} · ${voter.submittedAt}</p>
            </div>
            <div class="badge-row">
              <span class="badge ${voter.pvssVerified ? "badge-ok" : "badge-warn"}">PVSS ${voter.pvssVerified ? "OK" : "FAIL"}</span>
              <span class="badge ${voter.rangeVerified ? "badge-ok" : "badge-warn"}">ZKRP ${voter.rangeVerified ? "OK" : "FAIL"}</span>
              ${participantStatusBadge(voter.stake, chain.available ? "No Stake" : "Off-chain")}
            </div>
          </div>
          <p class="score-line">${formatScores(voter.scores, appState.snapshot.initiator.candidateLabels)}</p>
          <div class="token-list">
            ${voter.bindCommitments.map((item, index) => `<span class="token">Bind ${index + 1} · ${item}</span>`).join("")}
          </div>
          <div class="finance-panel">
            <div class="finance-strip">
              ${financeBadge("Required", chain.voterStakeEth)}
              ${financeBadge("Deposited", voter.stake.depositedEth)}
              ${financeBadge("Claimable", voter.stake.claimableEth)}
              ${financeBadge("Wallet", voter.stake.walletBalanceEth)}
            </div>
            <p class="record-meta account-line" title="${escapeHTML(voter.stake.address || "")}">${escapeHTML(accountLabel)}</p>
          </div>
          <div class="action-row">
            <button
              class="button-secondary"
              data-withdraw-voter-id="${voter.id}"
              ${!voter.stake.canWithdraw ? "disabled" : ""}
            >
              ${withdrawLabel}
            </button>
          </div>
          <details class="details">
            <summary>Show encrypted shares and ballot ciphertext</summary>
            <div class="details-grid">
              <div>${voter.encryptedShares.map((item, index) => `<span class="token token-soft">C${index + 1} · ${item}</span>`).join("")}</div>
              <div>${voter.ballotCipher.map((item, index) => `<span class="token token-soft">U${index + 1} · ${item}</span>`).join("")}</div>
            </div>
          </details>
        </article>
      `;
    })
    .join("");
}

function renderTalliers() {
  const { talliers, overview, chain } = appState.snapshot;
  const grid = document.getElementById("tallier-grid");

  grid.innerHTML = talliers
    .map((tallier) => {
      const mustStakeBeforeDecrypt = chain.available && !tallier.stake.staked;
      const canDecrypt = overview.canDecrypt && !tallier.hasDecrypted && !mustStakeBeforeDecrypt;
      const stakeLabel = !chain.available
        ? "Stake Unavailable"
        : tallier.stake.staked
          ? "Stake Locked"
          : `Stake ${formatETH(chain.tallierStakeEth)}`;
      const withdrawLabel = tallier.stake.withdrawn
        ? "Withdrawn"
        : tallier.stake.canWithdraw
          ? `Withdraw ${formatETH(tallier.stake.claimableEth)}`
          : "Withdraw After Settlement";
      const accountLabel = tallier.stake.address
        ? `Ganache account ${shortAddress(tallier.stake.address)}`
        : chain.available
          ? "Waiting for a tallier stake transaction."
          : "Ganache escrow unavailable.";

      return `
        <article class="tallier-card">
          <div class="record-head">
            <div>
              <h3>Tallier ${String(tallier.id).padStart(2, "0")}</h3>
              <p class="record-meta">${escapeHTML(tallier.publicKey)}</p>
            </div>
            <div class="badge-row">
              <span class="badge ${tallier.verified ? "badge-ok" : tallier.hasDecrypted ? "badge-warn" : "badge-neutral"}">
                ${tallier.hasDecrypted ? (tallier.verified ? "Verified Share" : "Share Published") : "Awaiting Share"}
              </span>
              ${participantStatusBadge(tallier.stake, chain.available ? "Unstaked" : "Off-chain")}
            </div>
          </div>
          <p class="tallier-copy">
            ${tallier.hasDecrypted ? `Share ${escapeHTML(tallier.share)}<br />Proof ${escapeHTML(tallier.proof)}` : chain.available ? "Stake first, then publish the verified decryption share for this tallier slot." : "Publish the verified decryption share without on-chain settlement because Ganache is unavailable."}
          </p>
          <div class="finance-panel">
            <div class="finance-strip">
              ${financeBadge("Required", chain.tallierStakeEth)}
              ${financeBadge("Deposited", tallier.stake.depositedEth)}
              ${financeBadge("Claimable", tallier.stake.claimableEth)}
              ${financeBadge("Wallet", tallier.stake.walletBalanceEth)}
            </div>
            <p class="record-meta account-line" title="${escapeHTML(tallier.stake.address || "")}">${escapeHTML(accountLabel)}</p>
          </div>
          <div class="tallier-actions">
            <button
              class="button-secondary"
              data-stake-id="${tallier.id}"
              ${!tallier.stake.canStake ? "disabled" : ""}
            >
              ${stakeLabel}
            </button>
            <button
              class="button-secondary"
              data-decrypt-id="${tallier.id}"
              ${!canDecrypt ? "disabled" : ""}
            >
              ${tallier.hasDecrypted ? "Decrypted" : "Decrypt Share"}
            </button>
            <button
              class="button-secondary"
              data-withdraw-tallier-id="${tallier.id}"
              ${!tallier.stake.canWithdraw ? "disabled" : ""}
            >
              ${withdrawLabel}
            </button>
          </div>
          <span class="record-meta">${tallier.decryptedAt || ""}</span>
        </article>
      `;
    })
    .join("");

  document.getElementById("finalize-button").disabled = !overview.canFinalize;
  document.getElementById("finalize-note").textContent = overview.canFinalize
    ? "Threshold satisfied. Finalizing will recover the tally and settle Ganache rewards."
    : `Need ${overview.threshold} verified talliers before final reconstruction and reward settlement.`;
}

function renderTally() {
  const { tally, chain, voters, talliers } = appState.snapshot;
  const board = document.getElementById("tally-board");

  if (!tally) {
    board.innerHTML = emptyState("The final tally will appear here once the decryption threshold is reached.");
    return;
  }

  const maxValue = Math.max(...tally.results.map((item) => item.total), 1);
  const rewardedVoters = voters.filter((voter) => voter.stake.honest).length;
  const rewardedTalliers = talliers.filter((tallier) => tallier.stake.honest).length;

  board.innerHTML = `
    <div class="badge-row">
      <span class="badge ${tally.verified ? "badge-ok" : "badge-warn"}">${tally.verified ? "Cross-check passed" : "Cross-check mismatch"}</span>
      <span class="badge ${chain.settled ? "badge-ok" : "badge-neutral"}">${chain.settled ? "Ganache settled" : "Ganache pending"}</span>
      <span class="record-meta">Finalized at ${tally.finalizedAt}</span>
    </div>
    <div class="reward-summary">
      ${rewardCard("Reward Pool", formatETH(chain.rewardPoolEth), chain.rewardSplit)}
      ${rewardCard("Initiator Claimable", formatETH(chain.initiator.claimableEth), chain.initiator.withdrawn ? "Already withdrawn" : `Account ${shortAddress(chain.initiator.address)}`)}
      ${rewardCard("Voter Rewards", `${rewardedVoters} honest voters`, "Claimable balances are shown on each voter record.")}
      ${rewardCard("Tallier Rewards", `${rewardedTalliers} honest talliers`, "Claimable balances are shown on each tallier card.")}
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
      <p class="metric-label">${escapeHTML(label)}</p>
      <p class="token-line">${escapeHTML(String(value))}</p>
    </article>
  `;
}

function ledgerCard(label, value, detail, badge = "") {
  return `
    <article class="stat-card ledger-card">
      <div class="card-topline">
        <p class="metric-label">${escapeHTML(label)}</p>
        ${badge}
      </div>
      <p class="ledger-value">${escapeHTML(String(value || "0"))}</p>
      <p class="record-meta">${escapeHTML(detail || "")}</p>
    </article>
  `;
}

function financeCard(label, finance, configuredStake, caption) {
  const address = finance.address || "";
  return `
    <article class="stat-card ledger-card finance-card">
      <div class="card-topline">
        <p class="metric-label">${escapeHTML(label)}</p>
        ${participantStatusBadge(finance, "Pending")}
      </div>
      <p class="ledger-value" title="${escapeHTML(address)}">${escapeHTML(shortAddress(address))}</p>
      <p class="record-meta">${escapeHTML(caption || "Role account")}</p>
      <div class="amount-grid">
        ${amountCell("Required", configuredStake)}
        ${amountCell("Deposited", finance.depositedEth)}
        ${amountCell("Claimable", finance.claimableEth)}
        ${amountCell("Wallet", finance.walletBalanceEth)}
      </div>
    </article>
  `;
}

function rewardCard(label, value, detail) {
  return `
    <article class="reward-card">
      <p class="metric-label">${escapeHTML(label)}</p>
      <p class="reward-value">${escapeHTML(value)}</p>
      <p class="record-meta">${escapeHTML(detail || "")}</p>
    </article>
  `;
}

function amountCell(label, value) {
  return `
    <span class="amount-cell">
      <span>${escapeHTML(label)}</span>
      <strong>${escapeHTML(formatETH(value))}</strong>
    </span>
  `;
}

function financeBadge(label, value) {
  return `<span class="token token-soft finance-token">${escapeHTML(label)} · ${escapeHTML(formatETH(value))}</span>`;
}

function statusBadge(label, tone = "neutral") {
  const className = tone === "ok" ? "badge-ok" : tone === "warn" ? "badge-warn" : "badge-neutral";
  return `<span class="badge ${className}">${escapeHTML(label)}</span>`;
}

function participantStatusBadge(finance, emptyLabel = "Pending") {
  if (!finance) {
    return statusBadge(emptyLabel, "neutral");
  }
  if (finance.withdrawn) {
    return statusBadge("Withdrawn", "neutral");
  }
  if (hasPositiveEth(finance.claimableEth)) {
    return statusBadge("Reward Ready", "ok");
  }
  if (finance.honest) {
    return statusBadge("Honest", "ok");
  }
  if (finance.staked) {
    return statusBadge("Staked", "ok");
  }
  return statusBadge(emptyLabel, "neutral");
}

function hasPositiveEth(value) {
  return Number(value || 0) > 0;
}

function formatETH(value) {
  const amount = value === undefined || value === null || value === "" ? "0.000" : String(value);
  return `${amount} ETH`;
}

function shortAddress(address) {
  if (!address || /^0x0{40}$/i.test(address)) {
    return "No account yet";
  }
  if (address.length <= 14) {
    return address;
  }
  return `${address.slice(0, 8)}...${address.slice(-6)}`;
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
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}
