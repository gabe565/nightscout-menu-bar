import { Browser, Call, Dialogs, Events, Window } from "/wails/runtime.js";
import { createFontPicker } from "./font-picker.js";

const service = "gabe565.com/nightscout-menu-bar/internal/settings.Service.";
const call = (method, ...args) => Call.ByName(service + method, ...args);

const $ = (id) => document.getElementById(id);
const form = $("form");
const errorMessage = (err) => err?.message ?? String(err);

const fontPicker = createFontPicker({
  input: $("font-input"),
  hidden: $("fontFile"),
  list: $("font-list"),
  onChange: () => form.dispatchEvent(new Event("input")),
});

// Form values <-> settings

function readForm() {
  const settings = {};
  for (const el of form.elements) {
    if (!el.name) continue;
    if (el.type === "checkbox") settings[el.name] = el.checked;
    else if (el.type === "number") settings[el.name] = Number(el.value);
    else settings[el.name] = el.value;
  }
  return settings;
}

function fillForm(settings) {
  for (const [name, value] of Object.entries(settings)) {
    const el = form.elements.namedItem(name);
    if (!el) continue;
    if (el.type === "checkbox") el.checked = value;
    else el.value = value;
  }
}

let saved = "";
const isDirty = () => JSON.stringify(readForm()) !== saved;
const updateDirty = () => ($("apply").disabled = !isDirty());

// Marks a field as saved without saving the rest of the form
function markSaved(name) {
  const settings = JSON.parse(saved);
  settings[name] = readForm()[name];
  saved = JSON.stringify(settings);
  updateDirty();
}

async function save() {
  $("save-error").textContent = "";
  const invalid = [...form.elements].find((el) => !validate(el));
  if (invalid) {
    const pane = invalid.closest('[role="tabpanel"]');
    selectTab(tabs.find((t) => t.dataset.pane === pane.id));
    invalid.focus();
    $("save-error").textContent = fieldName(invalid) + ": " + invalid.validationMessage;
    return false;
  }
  try {
    await call("Save", readForm());
    saved = JSON.stringify(readForm());
    $("apply").disabled = true;
    return true;
  } catch (err) {
    $("save-error").textContent = errorMessage(err);
    return false;
  }
}

// Tabs

const tabs = [...document.querySelectorAll('[role="tab"]')];

function selectTab(tab) {
  for (const t of tabs) {
    const selected = t === tab;
    t.setAttribute("aria-selected", selected);
    $(t.dataset.pane).hidden = !selected;
  }
  if (tab.dataset.pane === "display") updatePreview();
}

for (const tab of tabs) tab.addEventListener("click", () => selectTab(tab));

// Sizes the window to fit the tallest tab without scrolling
async function fitWindow() {
  const main = document.querySelector("main");
  const panes = [...document.querySelectorAll('[role="tabpanel"]')];
  const visible = panes.find((p) => !p.hidden);
  let tallest = 0;
  for (const pane of panes) {
    pane.hidden = false;
    for (const other of panes) if (other !== pane) other.hidden = true;
    tallest = Math.max(tallest, main.scrollHeight);
  }
  for (const pane of panes) pane.hidden = pane !== visible;

  const chrome = document.body.scrollHeight - main.clientHeight;
  await Window.SetSize(window.innerWidth, Math.ceil(chrome + tallest));
}

// Validation. The service validates again on save.

const durationPattern = /^(0|(\d+(\.\d+)?(ns|us|µs|ms|s|m|h))+)$/;

function validate(el) {
  if (!el.name || el.disabled) return true;
  const value = el.value.trim();
  let message = "";
  if ("duration" in el.dataset) {
    if (!durationPattern.test(value)) message = "Use a duration like 30s, 5m, or 1h30m.";
  } else if (el.name === "maxFontSize") {
    if (!(Number(value) > 0)) message = "Must be a positive number.";
  } else if (el.name === "url") {
    if (value && !/^https?:\/\/\S+$/i.test(value)) message = "Must start with http:// or https://.";
  }
  el.setCustomValidity(message);
  el.toggleAttribute("aria-invalid", message !== "");
  el.title = message;
  return message === "";
}

function fieldName(el) {
  return document.querySelector(`label[for="${el.id}"]`)?.textContent.trim() ?? el.name;
}

// Dependent fields

let platform = "";

function syncDisabled() {
  // Windows only shows an icon in the tray, so there's no text to preview
  $("preview-section").hidden = platform === "windows" && !$("dynamicIcon").checked;
  $("dynamic-icon-options").disabled = !$("dynamicIcon").checked;
  if (!$("dynamicIcon").checked) fontPicker.close();
  $("fontColor").disabled = $("autoColor").checked;
  $("fontColor-hex").textContent = $("autoColor").checked ? "" : $("fontColor").value;
}

// Connection test

async function testConnection() {
  const status = $("test-status");
  status.className = "status muted";
  status.textContent = "Testing…";
  $("test").disabled = true;
  try {
    const reading = await call("Test", readForm());
    status.className = "status ok";
    status.textContent = "✓ Connected. Latest reading: " + reading;
    updatePreview();
  } catch (err) {
    status.className = "status err";
    status.textContent = "✗ " + errorMessage(err);
  } finally {
    $("test").disabled = false;
  }
}

// Preview

let previewTimer;
let previewSeq = 0;

function schedulePreview() {
  clearTimeout(previewTimer);
  previewTimer = setTimeout(updatePreview, 150);
}

async function updatePreview() {
  const seq = ++previewSeq;
  const settings = readForm();
  const preview = await call("Preview", settings);
  if (seq !== previewSeq) return;

  const bar = $("preview-bar");
  const icon = $("preview-icon");
  const text = $("preview-text");
  const caption = $("preview-caption");

  bar.classList.remove("on-dark", "on-light");
  if (preview.icon) {
    // Use a background the icon color is visible on
    bar.classList.add(isLight(preview.iconColor) ? "on-dark" : "on-light");
    icon.src = preview.icon;
    icon.hidden = false;
    text.hidden = true;
  } else {
    icon.hidden = true;
    text.hidden = false;
    text.textContent = preview.text;
    text.classList.toggle("stale", preview.stale);
  }

  if (preview.error) {
    caption.className = "status err";
    caption.textContent = preview.error;
  } else {
    caption.className = "status muted";
    caption.textContent = preview.live ? "" : "Sample reading. Test the connection to preview your latest reading.";
  }
}

function isLight(hex) {
  const n = parseInt(hex.slice(1), 16);
  const [r, g, b] = [n >> 16, (n >> 8) & 0xff, n & 0xff];
  return 0.299 * r + 0.587 * g + 0.114 * b > 140;
}

// Nocturne

let connecting = null;

async function refreshNocturne() {
  const connected = !connecting && (await call("NocturneConnected", $("url").value));
  $("nocturne-status").textContent = connecting
    ? "Waiting for approval…"
    : connected ? "✓ Connected" : "Not connected";
  $("nocturne-button").textContent = connecting ? "Cancel" : connected ? "Disconnect" : "Connect…";
  $("nocturne-pending").hidden = !connecting;
  $("url").disabled = !!connecting;
  $("token-row").classList.toggle("dimmed", connected);
  $("token-inactive").hidden = !connected;
  return connected;
}

// The Nocturne option is only shown for servers that support it
let supportTimer;
let supportSeq = 0;

async function checkNocturneSupport() {
  const seq = ++supportSeq;
  const url = $("url").value.trim();
  const supported = !!connecting || (url !== "" && (await call("NocturneSupported", url)));
  if (seq !== supportSeq) return;
  for (const el of document.querySelectorAll("[data-nocturne]")) el.hidden = !supported;
}

let approvalURL = "";
Events.On("nocturne:code", (event) => {
  const { code, url } = event.data;
  $("nocturne-code").textContent = code;
  approvalURL = url;
});

$("nocturne-link").addEventListener("click", (e) => {
  e.preventDefault();
  if (approvalURL) Browser.OpenURL(approvalURL);
});

document.addEventListener("click", (e) => {
  const link = e.target.closest('a[href^="http"]');
  if (!link) return;
  e.preventDefault();
  Browser.OpenURL(link.href);
});

$("nocturne-button").addEventListener("click", async () => {
  if (connecting) {
    connecting.cancel();
    return;
  }
  if (await call("NocturneConnected", $("url").value)) {
    $("nocturne-status").textContent = "Disconnecting…";
    try {
      await call("DisconnectNocturne");
    } catch (err) {
      $("nocturne-status").textContent = errorMessage(err);
      return;
    }
    refreshNocturne();
    return;
  }

  $("nocturne-code").textContent = "…";
  approvalURL = "";
  connecting = call("ConnectNocturne", $("url").value);
  refreshNocturne();
  let connected = false;
  try {
    await connecting;
    connected = true;
  } catch (err) {
    if (err?.name !== "CancelError") {
      connecting = null;
      await refreshNocturne();
      $("nocturne-status").textContent = "✗ " + errorMessage(err);
      return;
    }
  }
  connecting = null;
  if (connected) {
    // The service saved the URL along with the connection
    markSaved("url");
    testConnection();
  }
  refreshNocturne();
});

// Events

form.addEventListener("input", (e) => {
  syncDisabled();
  updateDirty();
  if (validate(e.target) && $("save-error").textContent) $("save-error").textContent = "";
  if (e.target.id === "url") {
    refreshNocturne();
    clearTimeout(supportTimer);
    supportTimer = setTimeout(checkNocturneSupport, 600);
  }
  if (!$("display").hidden) schedulePreview();
});
form.addEventListener("change", () => {
  if (!$("display").hidden) schedulePreview();
});

// The service also routes the close button and Cmd+Q here
let closing = false;

async function requestClose() {
  if (closing) return;
  closing = true;
  try {
    if (isDirty()) {
      const choice = await Dialogs.Question({
        Title: "Do you want to save the changes to your settings?",
        Message: "Your changes will be lost if you don't save them.",
        Buttons: [
          { Label: "Save", IsDefault: true },
          { Label: "Don't Save" },
          { Label: "Cancel", IsCancel: true },
        ],
      });
      if (choice !== "Save" && choice !== "Don't Save") return;
      if (choice === "Save" && !(await save())) return;
    }
    await call("Close");
  } finally {
    closing = false;
  }
}

Events.On("settings:close-requested", requestClose);

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  if (await save()) call("Close");
});
$("apply").addEventListener("click", save);
$("cancel").addEventListener("click", requestClose);
$("test").addEventListener("click", testConnection);
$("browse-font").addEventListener("click", async () => {
  const path = await call("BrowseFont");
  if (path) {
    $("fontFile").value = path;
    fontPicker.sync();
    form.dispatchEvent(new Event("input"));
  }
});

document.addEventListener("keydown", (e) => {
  if (e.key === "Escape") {
    e.preventDefault();
    requestClose();
  } else if ((e.metaKey || e.ctrlKey) && e.key === "w") {
    e.preventDefault();
    requestClose();
  } else if ((e.metaKey || e.ctrlKey) && e.key >= "1" && e.key <= String(tabs.length)) {
    e.preventDefault();
    selectTab(tabs[Number(e.key) - 1]);
  }
});

// Startup

const [options, settings] = await Promise.all([call("Options"), call("Load")]);
platform = options.platform;
document.documentElement.dataset.platform = platform;
for (const [id, values] of [["units", options.units], ["logLevel", options.logLevels]]) {
  $(id).replaceChildren(...values.map((v) => new Option(v, v)));
}
fillForm(settings);
fontPicker.sync();
saved = JSON.stringify(readForm());
syncDisabled();
for (const el of form.elements) validate(el);
call("GuardClose");

$("reveal-token").addEventListener("click", () => {
  const token = $("token");
  const reveal = token.type === "password";
  token.type = reveal ? "text" : "password";
  $("reveal-token").setAttribute("aria-pressed", reveal);
  $("reveal-token").title = reveal ? "Hide token" : "Show token";
});
await refreshNocturne();
await fitWindow();
checkNocturneSupport();

// Slow or network-bound work after the window is ready
call("Fonts").then(fontPicker.setFonts);
if (settings.url) testConnection();
