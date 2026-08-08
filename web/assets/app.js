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
    el("p", {
      class: "tag",
      text: fam?.positioning || "We answer what should happen next.",
    }),
    el("h2", { class: "section-title", text: "Mission" }),
    el(
      "div",
      { class: "comps" },
      [
        el("div", { class: "comp" }, [
          el("div", { class: "k", text: "Not" }),
          el("div", { class: "meta", text: (fam?.mission_not || []).join(" · ") }),
        ]),
        el("div", { class: "comp" }, [
          el("div", { class: "k", text: "Instead" }),
          el("div", { class: "meta", text: (fam?.mission_instead || []).join(" · ") }),
        ]),
      ]
    ),
    el("p", {
      class: "meta",
      text: "Five questions: " + (fam?.five_questions || []).join(" → "),
    }),
    el("h2", { class: "section-title", text: "Innovation gate (ADR-9999)" }),
    el(
      "div",
      { class: "comps" },
      (fam?.innovation_gate || []).map((g) => el("div", { class: "comp" }, [el("div", { class: "k", text: g })]))
    ),
    el("h2", { class: "section-title", text: "Products + existence proofs" }),
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
          el("div", { class: "meta", text: "ADR-" + p.adr + (p.blueprint ? " · " + p.blueprint : "") }),
          el("div", { class: "meta", text: "Not: " + (p.commodity_avoid || []).join(", ") }),
          el(
            "ul",
            { class: "fam-features" },
            (p.features || []).map((f) =>
              el("li", {}, [
                el("strong", { text: f.name }),
                el("div", { class: "meta", text: f.novel_signal + " — " + (f.proof || "") }),
              ])
            )
          ),
          el("div", { class: "fam-scores", text: (p.novel_scores || []).slice(0, 4).join(" · ") }),
          el("div", {
            class: "meta",
            text: "Scores: loading…",
            "data-product": p.code,
          }),
        ])
      )
    ),
    el("div", {
      class: "meta",
      text: "Product APIs: GET /api/v1/products and /api/v1/products/{code}",
    }),
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
    el("h2", { class: "section-title", text: "Standards" }),
    el("p", { class: "meta", text: (fam?.standards || []).join(" · ") }),
    el("h2", { class: "section-title", text: "Research backlog" }),
    el("p", { class: "meta", text: (fam?.research_backlog || []).join(" · ") }),
    el("h2", { class: "section-title", text: "Maturity" }),
    el("p", { class: "meta", text: (fam?.maturity_levels || []).join(" → ") }),
  ]);
}

function duoView(impact, rec) {
  const layers = [
    ["business", impact?.business],
    ["engineering", impact?.engineering],
    ["operational", impact?.operational],
  ];
  return el("div", {}, [
    el("h2", { class: "section-title", text: "DUO — impact chain" }),
    el("p", { class: "meta", text: "Business → Engineering → Operational — not aggregated dashboards." }),
    el(
      "div",
      { class: "comps" },
      layers.map(([k, L]) =>
        el("div", { class: "comp" }, [
          el("div", { class: "k", text: L?.label || k }),
          el("div", { class: "v", text: L?.score != null ? String(L.score) : "—" }),
          el("div", { class: "meta", text: L?.summary || "" }),
        ])
      )
    ),
    el("h2", { class: "section-title", text: "Recommended action" }),
    rec
      ? el("div", { class: "src" }, [
          el("div", {}, [
            el("strong", { text: rec.title || "Action" }),
            el("div", { class: "meta", text: rec.action || "" }),
            el("div", {
              class: "meta",
              text:
                "confidence " +
                (rec.confidence != null ? Math.round(rec.confidence * 100) + "%" : "—") +
                " · effort " +
                (rec.estimated_effort || "—") +
                " · " +
                (rec.expected_impact || ""),
            }),
            el("div", { class: "fam-scores", text: (rec.evidence || []).join(" · ") }),
          ]),
        ])
      : el("p", { class: "meta", text: "No recommendation yet." }),
    el("h2", { class: "section-title", text: "Source scores" }),
    el(
      "div",
      { class: "sources" },
      (impact?.sources || []).map((s) =>
        el("div", { class: "src" }, [
          el("strong", { text: s.product + "." + s.name }),
          el("span", {
            class: "pill " + (s.mode === "live" ? "ok" : s.mode === "demo" ? "muted" : "bad"),
            text: s.mode + " " + s.value,
          }),
        ])
      )
    ),
  ]);
}

function campaignsView(campaigns, channels) {
  const list = campaigns || [];
  const chans = channels || [];
  return el("div", {}, [
    el("h2", { class: "section-title", text: "Campaigns (DPO F3)" }),
    el("p", {
      class: "meta",
      text: "Presence engineering — dry-run channel artifacts, then arm/send when credentials exist (ADR-0402).",
    }),
    el("h2", { class: "section-title", text: "Channels" }),
    el(
      "div",
      { class: "sources" },
      chans.map((c) =>
        el("div", { class: "src" }, [
          el("div", {}, [
            el("strong", { text: c.name || c.id }),
            el("div", { class: "meta", text: c.kind + " · limit " + c.char_limit + (c.setup_hint ? " — " + c.setup_hint : "") }),
          ]),
          el("span", {
            class: "pill " + (c.live_status === "ready" ? "ok" : "muted"),
            text: c.live_status,
          }),
        ])
      )
    ),
    el("h2", { class: "section-title", text: "Campaigns" }),
    el(
      "div",
      { class: "sources" },
      list.length
        ? list.map((camp) =>
            el("div", { class: "src", style: "flex-direction:column;align-items:stretch;gap:0.5rem" }, [
              el("div", { style: "display:flex;justify-content:space-between;gap:1rem;flex-wrap:wrap" }, [
                el("div", {}, [
                  el("strong", { text: camp.id }),
                  el("div", { class: "meta", text: camp.title + " · " + camp.status }),
                ]),
                el("div", { class: "actions", style: "margin:0" }, [
                  el("button", {
                    text: "Dry-run render",
                    onclick: async () => {
                      await fetch("/api/v1/campaigns/" + camp.id + "/render", { method: "POST" });
                      boot("campaigns");
                    },
                  }),
                  el("button", {
                    text: "Arm",
                    onclick: async () => {
                      await fetch("/api/v1/campaigns/" + camp.id + "/arm", { method: "POST" });
                      boot("campaigns");
                    },
                  }),
                  el("button", {
                    text: "Send / manual receipts",
                    onclick: async () => {
                      await fetch("/api/v1/campaigns/" + camp.id + "/send", { method: "POST" });
                      boot("campaigns");
                    },
                  }),
                ]),
              ]),
              el(
                "div",
                { class: "comps" },
                (camp.channels || [])
                  .filter((ch) => (ch.artifacts || []).length)
                  .map((ch) => {
                    const art = ch.artifacts[0];
                    return el("div", { class: "comp" }, [
                      el("div", { class: "k", text: ch.channel_id + " (" + ch.status + ")" }),
                      el("div", {
                        class: "meta",
                        text:
                          (art.subject || art.title || "") +
                          (art.over_limit ? " ⚠ over limit" : "") +
                          " · " +
                          art.char_count +
                          "/" +
                          art.char_limit,
                      }),
                      el("pre", {
                        class: "meta",
                        style: "white-space:pre-wrap;max-height:8rem;overflow:auto",
                        text: art.body,
                      }),
                      el("button", {
                        text: "Copy",
                        onclick: () => navigator.clipboard.writeText(art.body),
                      }),
                    ]);
                  })
              ),
            ])
          )
        : [el("p", { class: "meta", text: "No campaigns seeded." })]
    ),
  ]);
}

function render(tab, score, sources, meta, eng, content, baselines, fam, impact, rec, campaigns, channels) {
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
      class: tab === "duo" ? "nav-btn active" : "nav-btn",
      text: "DUO",
      onclick: () => boot("duo"),
    }),
    el("button", {
      class: tab === "campaigns" ? "nav-btn active" : "nav-btn",
      text: "Campaigns",
      onclick: () => boot("campaigns"),
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
  } else if (tab === "duo") {
    body = duoView(impact, rec);
  } else if (tab === "campaigns") {
    body = campaignsView(campaigns, channels);
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
          const label = prompt("Baseline label", "pre-launch");
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
  const [score, sources, meta, eng, content, baselines, fam, impact, rec, productList, campaigns, channels] =
    await Promise.all([
      getJSON("/api/v1/score"),
      getJSON("/api/v1/sources/status"),
      getJSON("/api/v1/meta"),
      getJSON("/api/v1/engineering"),
      getJSON("/api/v1/content"),
      getJSON("/api/v1/baselines"),
      getJSON("/api/v1/family"),
      getJSON("/api/v1/duo/impact"),
      getJSON("/api/v1/duo/recommend"),
      getJSON("/api/v1/products").catch(() => []),
      getJSON("/api/v1/campaigns").catch(() => []),
      getJSON("/api/v1/channels").catch(() => []),
    ]);
  render(tab, score, sources, meta, eng, content, baselines, fam, impact, rec, campaigns, channels);
  if (tab === "family" && Array.isArray(productList)) {
    const by = Object.fromEntries(productList.map((p) => [p.code, p]));
    document.querySelectorAll("[data-product]").forEach((node) => {
      const p = by[node.dataset.product];
      if (!p) return;
      const bits = (p.scores || []).map((s) => s.name + "=" + s.value + "(" + s.mode + ")");
      node.textContent = bits.length ? bits.join(" · ") : "no scores yet — run collectors";
    });
  }
}

boot().catch((e) => {
  document.getElementById("app").textContent = "Failed to load Observatory Platform: " + e.message;
});
