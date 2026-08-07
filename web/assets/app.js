async function getJSON(path) {
  const r = await fetch(path);
  if (!r.ok) throw new Error(`${path} ${r.status}`);
  return r.json();
}

function el(tag, attrs = {}, kids = []) {
  const n = document.createElement(tag);
  Object.entries(attrs).forEach(([k, v]) => {
    if (k === "class") n.className = v;
    else if (k === "text") n.textContent = v;
    else if (k.startsWith("on") && typeof v === "function") n.addEventListener(k.slice(2).toLowerCase(), v);
    else n.setAttribute(k, v);
  });
  kids.forEach((c) => n.append(c));
  return n;
}

function render(score, sources, meta) {
  const root = document.getElementById("app");
  root.innerHTML = "";
  const value = score?.value ?? 0;
  const comps = score?.components || {};

  const shell = el("div", { class: "shell" }, [
    el("h1", { class: "brand", text: "Digital Presence Observatory" }),
    el("p", {
      class: "tag",
      text: "Engineering observability for digital presence — Observe. Measure. Improve. Pilot tenant " + (meta?.tenant || "dasmlab.org") + ".",
    }),
    el("div", { class: "hero-score" }, [
      (() => {
        const ring = el("div", { class: "ring" });
        ring.style.setProperty("--p", String(value));
        ring.append(
          el("div", {}, [
            el("div", { class: "val", text: String(value) }),
            el("div", { class: "lbl", text: "Overall" }),
          ])
        );
        return ring;
      })(),
      el(
        "div",
        { class: "comps" },
        Object.entries(comps).map(([k, v]) =>
          el("div", { class: "comp" }, [
            el("div", { class: "k", text: k }),
            el("div", { class: "v", text: String(v) }),
          ])
        )
      ),
    ]),
    el("h2", { class: "section-title", text: "Collector health" }),
    el(
      "div",
      { class: "sources" },
      (sources || []).map((s) =>
        el("div", { class: "src" }, [
          el("div", {}, [
            el("strong", { text: s.name }),
            el("div", { class: "meta", text: s.message || s.last_error || "" }),
          ]),
          el("span", { class: "pill " + (s.healthy ? "ok" : "bad"), text: s.healthy ? "Healthy" : "Failing" }),
        ])
      )
    ),
    el("div", { class: "actions" }, [
      el("button", {
        text: "Run collectors",
        onclick: async () => {
          await fetch("/api/v1/collect/run", { method: "POST" });
          setTimeout(boot, 800);
        },
      }),
      el("button", {
        text: "Refresh",
        onclick: () => boot(),
      }),
    ]),
    el("p", {
      class: "meta",
      text: `DOP · DPO ${meta?.version || ""} · Five questions: ${(meta?.five_questions || []).join(" → ")}`,
    }),
  ]);
  root.append(shell);
}

async function boot() {
  const [score, sources, meta] = await Promise.all([
    getJSON("/api/v1/score"),
    getJSON("/api/v1/sources/status"),
    getJSON("/api/v1/meta"),
  ]);
  render(score, sources, meta);
}

boot().catch((e) => {
  document.getElementById("app").textContent = "Failed to load DPO: " + e.message;
});
