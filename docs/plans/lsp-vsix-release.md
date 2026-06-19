# Plan: Release the VSIX with the Rust LSP (start to finish)

Ship the Osprey VS Code extension carrying the new Rust language server, on the
Nimblesite OIDC publishing path — then extend the same `osprey lsp` engine to
Open VSX, Neovim, and Zed.

**Spec:** [0020-LanguageServerAndEditors.md](../specs/0020-LanguageServerAndEditors.md).
**Standards:** [Shipwright](https://github.com/Nimblesite/Shipwright) (public
recipe) · [Nimblesite/NimblesiteDeployment](https://github.com/Nimblesite/NimblesiteDeployment)
(private concrete infra) · [lspkit](https://github.com/Nimblesite/lspkit).

## Where we are

- ✅ Rust LSP merged ([#137](https://github.com/MelbourneDeveloper/osprey/pull/137)):
  `crates/osprey-lsp` on published `lspkit-*` crates; TypeScript server deleted;
  VS Code client spawns `osprey lsp`; webcompiler bridge spawns the same binary.
- ✅ `shipwright.json` declares `osprey-compiler` (CLI = the LSP) + `osprey-vscode`.
- ✅ `release.yml` `vsix` job builds a **per-platform VSIX** with the
  version-matched binary bundled and verified inside the package.
- ⚠️ **Blocker:** that job publishes with a legacy **`VSCE_PAT`** —
  the exact stored-token path NimblesiteDeployment is retiring in favour of OIDC.
- ⚠️ **Blocker:** osprey is `MelbourneDeveloper/osprey`. The shared Entra app's
  federated credential is a **wildcard for `repo:Nimblesite/*:environment:release`
  only** — it does **not** trust `MelbourneDeveloper/*`. OIDC will not work until
  this is resolved.
- ⬜ No Open VSX, Neovim, or Zed integration yet.

## The linchpin: repo-owner mismatch

The OIDC flow trusts `repo:Nimblesite/*`. Osprey is under `MelbourneDeveloper`.

**Decision: Option A.** Osprey is years from production-useful; the rest of the
Nimblesite shelf (Shipwright, lspkit, SharpLsp, Deslop) is usable today. Moving a
not-ready language into the org dilutes that "use this now" signal, and the move is
roughly one-way (URLs, tap/bucket refs, CI) — paying that cost to advertise
something unfinished is the wrong trade. A dedicated federated credential gets the
release out with one trivial, reversible Azure object and keeps Osprey at arm's
length until it earns the org namespace.

| Option | What | Trade-off | Verdict |
|---|---|---|---|
| **A. Dedicated federated credential** | Add a standard subject credential `repo:MelbourneDeveloper/osprey:environment:release` to the shared app. | One Azure object, scoped to exactly this repo. No migration. | **✅ Chosen.** |
| B. Second wildcard | Add `repo:MelbourneDeveloper/*:environment:release`. | Trusts every future MelbourneDeveloper repo — broader than needed. | Only if more such repos are coming. |
| C. Transfer repo to Nimblesite | Move osprey under the org; wildcard covers it. | Covered by existing trust, but a large, one-way move (URLs, tap repo refs, CI), and it puts an unfinished product on the org shelf. | **Deferred** — see trigger below. |

**Option C trigger (revisit later):** transfer osprey into the `Nimblesite` org
once it is production-useful — a tagged `v1.0.0`-class release that a real user
could adopt. At that point org membership signals the right thing and the wildcard
covers it for free; until then it stays under `MelbourneDeveloper`.

Osprey already publishes as `nimblesite.osprey` (the **same** publisher), so the
shared app is already a Contributor member — **no manual Marketplace member step
is needed.** Only the federated credential is missing.

**Mark the listing as a preview.** While Osprey is pre-production, set
`"preview": true` in `vscode-extension/package.json` so the Marketplace badge is
honest about maturity. Drop it at the same `v1.0.0` milestone that triggers the
Option C move.

---

## Phase 1 — Publish the VSIX via OIDC (the immediate goal)

Goal: tag a release → per-platform VSIX with the Rust LSP lands on the Marketplace
with **no PAT**.

### 1.1 Azure (one-time, via NimblesiteDeployment) — Option A

Run by someone with `az` signed into the Nimblesite tenant. From
[`scripts/setup-marketplace-oidc.sh`](https://github.com/Nimblesite/NimblesiteDeployment)
conventions:

```sh
APP_OBJECT_ID=d2586e67-f2dc-497d-97a4-2edcb885ec96   # shared app object id (azure-inventory.md)
az ad app federated-credential create --id "$APP_OBJECT_ID" --parameters '{
  "name": "github-osprey-release",
  "issuer": "https://token.actions.githubusercontent.com",
  "subject": "repo:MelbourneDeveloper/osprey:environment:release",
  "audiences": ["api://AzureADTokenExchange"]
}'
```

- [ ] Federated credential created for `MelbourneDeveloper/osprey`.

### 1.2 GitHub config on `MelbourneDeveloper/osprey`

```sh
gh api -X PUT repos/MelbourneDeveloper/osprey/environments/release

gh secret set AZURE_CLIENT_ID --env release --repo MelbourneDeveloper/osprey \
  --body "beacf14a-c783-4bab-80a6-dd4936cb1da3"   # shared app client id
gh secret set AZURE_TENANT_ID --env release --repo MelbourneDeveloper/osprey \
  --body "0a282151-85df-4a81-b083-52221a26d8e7"   # tenant id
```

- [ ] `release` environment exists with `AZURE_CLIENT_ID` + `AZURE_TENANT_ID`.

### 1.3 Convert the `vsix` job to OIDC (`release.yml`)

Replace the `VSCE_PAT` publish step. The job must bind to the `release`
environment and mint a Marketplace token from the OIDC session:

```yaml
  vsix:
    name: Publish VSIX ${{ matrix.target }}
    needs: [version, build]
    runs-on: ${{ matrix.os }}
    environment: release            # <-- gates OIDC trust
    permissions:
      id-token: write               # <-- request the OIDC token
      contents: read
    # ... existing matrix, checkout, build/stage/package steps unchanged ...
      - name: Azure login (OIDC, no secret)
        uses: azure/login@v2
        with:
          client-id: ${{ secrets.AZURE_CLIENT_ID }}
          tenant-id: ${{ secrets.AZURE_TENANT_ID }}
          allow-no-subscriptions: true
      - name: Publish to Marketplace (OIDC)
        if: ${{ vars.SKIP_VSCE_PUBLISH != 'true' }}
        working-directory: vscode-extension
        run: |
          set -euo pipefail
          TOKEN=$(az account get-access-token \
            --resource 499b84ac-1321-427f-aa17-267ca6975798 \
            --query accessToken -o tsv)        # immutable vsce/Azure DevOps resource id
          npx vsce publish --packagePath "$(ls *.vsix)" --pat "$TOKEN"
```

- [ ] `environment: release` + `id-token: write` added to the `vsix` job.
- [ ] `azure/login@v2` + token-mint + `vsce publish --pat $TOKEN` replace the
      `VSCE_PAT` step.
- [ ] `VSCE_PAT` removed from the workflow (and from repo secrets once green).
- [ ] Keep the `SKIP_VSCE_PUBLISH` dry-run escape hatch.
- [ ] Cross-check the job shape against Shipwright's
      `templates/gh-actions/publish-vsix-per-platform.yml`.

### 1.4 Verify (pre-tag)

- [ ] `make ci` green locally.
- [ ] Dry-run: set repo var `SKIP_VSCE_PUBLISH=true`, push a throwaway tag, confirm
      VSIX builds, the bundled `osprey` binary is verified inside each `.vsix`, and
      OIDC login succeeds (the publish step is skipped).
- [ ] Confirm `npm run test:shipwright` passes (manifest validation in the
      `version` job).

---

## Phase 2 — Open VSX `[EDITOR-OPENVSX]`

Same per-platform VSIX, second registry — for VSCodium/Cursor/Windsurf/Gitpod.
Open VSX has **no** OIDC; it uses a long-lived `OVSX_PAT`, kept independent.

- [ ] Create an Open VSX `nimblesite` namespace + token; store `OVSX_PAT` as a repo
      secret.
- [ ] Add a `publish-openvsx` matrix job (mirrors `vsix`, **no** `environment`/OIDC):
      `npx ovsx publish "$(ls *.vsix)" -p "$OVSX_PAT" --target ${{ matrix.target }}`.
- [ ] Make it independent of the Marketplace job (neither gates the other);
      `deploy-site`/`deploy-webcompiler` should not block on Open VSX.
- [ ] Document the install path in `vscode-extension/README.md`.

---

## Phase 3 — Onboard Osprey into NimblesiteDeployment

The private ops repo is the source of truth for "what's wired". A PR/commit there:

- [ ] Add a `github-osprey-release` row to the **federated identity credentials**
      table in `docs/azure-inventory.md` (note: subject credential, not the
      wildcard, because osprey is outside `Nimblesite/*`).
- [ ] Add `MelbourneDeveloper/osprey` → `nimblesite.osprey` to the **Onboarded
      repositories** table (Marketplace OIDC ✅; Open VSX ✅ once Phase 2 ships).
- [ ] In `docs/onboarding-a-new-vsix.md`, add a short note that **non-`Nimblesite/*`
      repos need a dedicated federated credential** (the wildcard doesn't cover
      them) — osprey is the first such case.
- [ ] Update `docs/RELEASING.md` here: drop `VSCE_PAT` from the secrets table; add
      `AZURE_CLIENT_ID`/`AZURE_TENANT_ID` (env: release) and `OVSX_PAT`; note the
      OIDC publishing model.

---

## Phase 4 — Shipwright conformance

- [ ] `shipwright.json` unchanged for the LSP: it is the `osprey` CLI via the `lsp`
      subcommand, **not** a separate component — nothing to add, nothing to drift.
- [ ] Confirm the version contract end-to-end: `osprey --version` →
      `osprey X.Y.Z`; `osprey --version --json` matches the version manifest schema
      (`[SWR-VERSION-CLI-OUTPUT]`).
- [ ] Confirm no source field leaves `0.0.0-dev` (`[SWR-VERSION-BUILD-STAMPING]`):
      `Cargo.toml`, `vscode-extension/package.json`, `shipwright.json`.
- [ ] Keep the "bundled binary present inside the VSIX" assertion (acceptance gate).
- [ ] Activation verifies the bundled compiler against the manifest
      (`hosts.vscode.onMismatch: prompt-reinstall`) via `@nimblesite/shipwright-vscode`.

---

## Phase 5 — Neovim `[EDITOR-NEOVIM]` (next)

The server is editor-agnostic; only a client recipe is missing.

- [ ] Add an `editors/neovim/` recipe + README: `vim.lsp` / `nvim-lspconfig`
      registration pointing `cmd = { 'osprey', 'lsp' }` at the PATH binary
      (snippet already in the spec).
- [ ] Document install via Homebrew/Scoop/GitHub release (the binary is the server).
- [ ] (Optional) Upstream a server definition to `nvim-lspconfig`.

## Phase 6 — Zed `[EDITOR-ZED]` (next)

- [ ] Create a Zed extension (Rust→WASM) registering Osprey + its language server
      command `osprey lsp`, using `shipwright-zed` for version-matching/acquisition.
- [ ] Publish to the Zed extension registry.

## Phase 7 — Cut the release

- [ ] `git tag vX.Y.Z && git push origin vX.Y.Z` from an up-to-date `main`.
- [ ] Watch `release.yml`: build → GitHub Release → brew → scoop → **VSIX (OIDC)**
      → Open VSX → site/webcompiler.
- [ ] Smoke test: install from Marketplace + Open VSX on each platform; open a
      `.osp` file; confirm diagnostics, hover, go-to-definition, signature help; and
      that the bundled `osprey lsp` version matches the extension version.

---

## Risks & decisions

- **Repo-owner mismatch (highest):** OIDC fails until Phase 1.1. Decide A/B/C above
  first; A is recommended.
- **`VSCE_PAT` removal:** delete the secret only after an OIDC publish is confirmed
  green, so there's a rollback path during the transition.
- **`OVSX_PAT` is the one standing token** in the publishing path — scope it to the
  `nimblesite` Open VSX namespace only; it is never shared with the Marketplace
  identity.
- **Version drift is structurally impossible** for the LSP: the extension bundles
  the exact binary it launches, and that binary *is* the server. Keep it that way —
  do not introduce a separately-versioned server artifact.
