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

function metricGrid(title, obj) {
  const entries = Object.entries(obj || {}).filter(([, v]) => v !== undefined && v !== null);
  return el("div", {}, [
    el("h2", { class: "section-title", text: title }),
    el(
      "div",
      { class: "comps" },
      entries.length
        ? entries.map(([k, v]) =>
            el("div", { class: "comp" }, [
              el("div", { class: "k", text: k }),
              el("div", { class: "v", text: typeof v === "number" ? String(Math.round(v * 10) / 10) : String(v) }),
            ])
          )
        : [el("p", { class: "meta", text: "No samples yet — run collectors." })]
    ),
  ]);
}

function contentTable(paths) {
  const rows = paths || [];
  return el("div", {}, [
    el("h2", { class: "section-title", text: "Content spine" }),
    rows.length
      ? el(
          "div",
          { class: "sources" },
          rows.map((p) =>
            el("div", { class: "src" }, [
              el("div", {}, [
                el("strong", { text: p.path }),
                el("div", {
                  class: "meta",
                  text: `views ${p.page_views || 0} · gsc ${p.impressions || 0}/${p.clicks || 0} · bots ${p.bot_hits || 0}`,
                }),
              ]),
            ])
          )
        )
      : el("p", { class: "meta", text: "No path entities yet — need live sitemap + collect." }),
  ]);
}

function familyView(fam) {
  const products = fam?.products || [];
  return el("div", {}, [
    el("h2", { class: "section-title", text: "Observatory family" }),
    el("p", {
      class: "meta",
      text: fam?.differentiator || "Discover things nobody measures today.",
    }),
    el("h2", { class: "section-title", text: "Innovation gate" }),
    el(
      "div",
      { class: "comps" },
      (fam?.innovation_gate || []).map((g) => el("div", { class: "comp" }, [el("div", { class: "k", text: g })]))
    ),
    el("h2", { class: "section-title", text: "Products" }),
    el(
      "div",
      { class: "family-grid" },
      products.map((p) =>
        el("div", { class: "fam-card" + (p.status === "live" ? " live" : "") }, [
          el("div", { class: "fam-head" }, [
            el("strong", { text: p.code.toUpperCase() }),
            el("span", { class: "pill " + (p.status === "live" ? "ok" : "muted"), text: p.status }),
          ]),
          el("div", { class: "fam-name", text: p.name }),
          el("div", { class: "meta", text: "ADR-" + p.adr }),
          el("div", { class: "meta", text: "Not: " + (p.commodity_avoid || []).join(", ") }),
          el("div", { class: "fam-scores", text: (p.novel_scores || []).join(" · ") }),
        ])
      )
    ),
    el("h2", { class: "section-title", text: "Shared pipeline" }),
    el("p", { class: "meta", text: (fam?.architecture_pipeline || []).join(" → ") }),
    el("h2", { class: "section-title", text: "Platform layers" }),
    el(
      "div",
      { class: "sources" },
      (fam?.platform_layers || []).map((l) =>
        el("div", { class: "src" }, [
          el("strong", { text: l.name }),
          el("div", { class: "meta", text: (l.elements || []).join(" · ") }),
        ])
      )
    ),
    el("h2", { class: "section-title", text: "Shared SDKs" }),
    el("p", { class: "meta", text: (fam?.shared_sdks || []).join(" · ") }),
    el("h2", { class: "section-title", text: "Research backlog" }),
    el("p", { class: "meta", text: (fam?.research_backlog || []).join(" · ") }),
    el("h2", { class: "section-title", text: "Maturity" }),
    el("p", { class: "meta", text: (fam?.maturity_levels || []).join(" → ") }),
  ]);
}

function render(tab, score, sources, meta, eng, content, baselines, fam) {
  const root = document.getElementById("app");
  root.innerHTML = "";
  const value = score?.value ?? 0;
  const comps = score?.components || {};

  const nav = el("div", { class: "nav" }, [
    el("button", {
      class: tab === "family" ? "nav-btn active" : "nav-btn",
      text: "Family",
      onclick: () => boot("family"),
    }),
    el("button", {
      class: tab === "executive" ? "nav-btn active" : "nav-btn",
      text: "DPO Executive",
      onclick: () => boot("executive"),
    }),
    el("button", {
      class: tab === "engineering" ? "nav-btn active" : "nav-btn",
      text: "DPO Engineering",
      onclick: () => boot("engineering"),
    }),
  ]);

  let body;
  if (tab === "family") {
    body = familyView(fam);
  } else if (tab === "engineering") {
    body = el("div", {}, [
      metricGrid("Bot / crawl", eng?.bots),
      metricGrid("Index / technical", eng?.index),
      metricGrid("Search (GSC)", eng?.search),
      metricGrid("GitHub authority", eng?.authority),
      metricGrid("Engagement (Activity)", eng?.engagement),
      contentTable(content?.paths),
      el("h2", { class: "section-title", text: "Collector health" }),
      el(
        "div",
        { class: "sources" },
        (sources || []).map((s) =>
          el("div", { class: "src" }, [
            el("div", {}, [
              el("strong", { text: s.name }),
              el("div", { class: "meta", text: (s.message || "") + (s.last_error ? " — " + s.last_error : "") }),
            ]),
            el("span", { class: "pill " + (s.healthy ? "ok" : "bad"), text: s.healthy ? "Healthy" : "Failing" }),
          ])
        )
      ),
    ]);
  } else {
    body = el("div", {}, [
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
      el("h2", { class: "section-title", text: "Baselines" }),
      el(
        "div",
        { class: "sources" },
        (baselines || []).length
          ? baselines.map((b) =>
              el("div", { class: "src" }, [
                el("strong", { text: b.label }),
                el("span", { class: "meta", text: b.created_at }),
              ])
            )
          : [el("p", { class: "meta", text: "No baselines yet — freeze one after collect." })]
      ),
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
    ]);
  }

  const shell = el("div", { class: "shell" }, [
    el("h1", { class: "brand", text: "Observatory Platform" }),
    el("p", {
      class: "tag",
      text:
        (fam?.tagline || "Engineering observability for complex systems.") +
        " Active product: " +
        (meta?.product || "dpo").toUpperCase() +
        " · pilot " +
        (meta?.tenant || "dasmlab.org") +
        ".",
    }),
    nav,
    body,
    el("div", { class: "actions" }, [
      el("button", {
        text: "Run collectors",
        onclick: async () => {
          await fetch("/api/v1/collect/run", { method: "POST" });
          setTimeout(() => boot(tab), 1500);
        },
      }),
      el("button", {
        text: "Freeze baseline",
        onclick: async () => {
          const label = prompt("Baseline label", "pre-home-2.0");
          if (!label) return;
          await fetch("/api/v1/baseline", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ label }),
          });
          boot(tab);
        },
      }),
      el("button", { text: "Refresh", onclick: () => boot(tab) }),
    ]),
    el("p", {
      class: "meta",
      text: `DOP · ${(meta?.product || "dpo").toUpperCase()} ${meta?.version || ""} · ${(meta?.five_questions || []).join(" → ")}`,
    }),
  ]);
  root.append(shell);
}

async function boot(tab = "family") {
  const [score, sources, meta, eng, content, baselines, fam] = await Promise.all([
    getJSON("/api/v1/score"),
    getJSON("/api/v1/sources/status"),
    getJSON("/api/v1/meta"),
    getJSON("/api/v1/engineering"),
    getJSON("/api/v1/content"),
    getJSON("/api/v1/baselines"),
    getJSON("/api/v1/family"),
  ]);
  render(tab, score, sources, meta, eng, content, baselines, fam);
}

boot().catch((e) => {
  document.getElementById("app").textContent = "Failed to load Observatory Platform: " + e.message;
});
