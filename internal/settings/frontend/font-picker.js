// A font dropdown that shows each font's name in the font itself.
// The selected font is saved in a hidden input, since the visible input shows its name.

const families = new Map();

// Loads a font for previewing and returns its CSS font family
function previewFamily(file) {
  if (!families.has(file)) {
    const family = `font-preview-${families.size}`;
    const face = new FontFace(family, `url("/font?file=${encodeURIComponent(file)}")`);
    document.fonts.add(face);
    face.load().catch(() => {});
    families.set(file, family);
  }
  return `"${families.get(file)}", system-ui`;
}

export function createFontPicker({ input, hidden, list, onChange }) {
  const embedded = { file: "", name: "Embedded font" };
  let fonts = [embedded];
  let shown = [];
  let active = -1;

  // Fonts only load once they scroll into view
  const observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (!entry.isIntersecting) continue;
        entry.target.firstChild.style.fontFamily = previewFamily(entry.target.dataset.file);
        observer.unobserve(entry.target);
      }
    },
    { root: list },
  );

  const nameOf = (file) => fonts.find((f) => f.file === file)?.name ?? file;

  function sync() {
    input.value = hidden.value === "" ? "" : nameOf(hidden.value);
    input.style.fontFamily = previewFamily(hidden.value);
  }

  function select(file) {
    close();
    if (file !== hidden.value) {
      hidden.value = file;
      onChange();
    }
    sync();
  }

  // Accepts a typed font name, file name, or path
  function commit() {
    const text = input.value.trim();
    if (text === nameOf(hidden.value) || (text === "" && hidden.value === "")) return;
    const match = fonts.find((f) => f.name.toLowerCase() === text.toLowerCase());
    select(match ? match.file : text);
  }

  function render(filter) {
    const query = filter.trim().toLowerCase();
    shown = fonts.filter(
      (f) => !query || f.name.toLowerCase().includes(query) || f.file.toLowerCase().includes(query),
    );
    list.replaceChildren(
      ...shown.map((font) => {
        const option = document.createElement("div");
        option.className = "combo-option";
        option.setAttribute("role", "option");
        option.dataset.file = font.file;
        const name = document.createElement("span");
        name.className = "name";
        name.textContent = font.name;
        const file = document.createElement("span");
        file.className = "file";
        file.textContent = font.file;
        option.append(name, file);
        // Select before the input loses focus
        option.addEventListener("mousedown", (e) => {
          e.preventDefault();
          select(font.file);
        });
        observer.observe(option);
        return option;
      }),
    );
    setActive(shown.findIndex((f) => f.file === hidden.value));
  }

  function setActive(index) {
    active = index;
    [...list.children].forEach((option, i) => {
      option.classList.toggle("active", i === index);
      option.setAttribute("aria-selected", i === index);
    });
    list.children[index]?.scrollIntoView({ block: "nearest" });
  }

  function open() {
    if (input.disabled || !list.hidden) return;
    // Fixed position so the list isn't clipped by the scrolling page
    const rect = input.getBoundingClientRect();
    const below = window.innerHeight - rect.bottom - 8;
    list.style.left = `${rect.left}px`;
    list.style.width = `${rect.width}px`;
    if (below >= 160 || below >= rect.top) {
      list.style.top = `${rect.bottom + 4}px`;
      list.style.bottom = "";
      list.style.maxHeight = `${Math.min(260, below)}px`;
    } else {
      list.style.top = "";
      list.style.bottom = `${window.innerHeight - rect.top + 4}px`;
      list.style.maxHeight = `${Math.min(260, rect.top - 8)}px`;
    }
    list.hidden = false;
    input.setAttribute("aria-expanded", "true");
    render("");
  }

  function close() {
    list.hidden = true;
    input.setAttribute("aria-expanded", "false");
  }

  input.addEventListener("focus", open);
  input.addEventListener("click", open);
  input.addEventListener("input", () => {
    open();
    render(input.value);
  });
  input.addEventListener("blur", () => {
    close();
    commit();
  });
  input.addEventListener("keydown", (e) => {
    if (e.key === "ArrowDown" || e.key === "ArrowUp") {
      e.preventDefault();
      open();
      const step = e.key === "ArrowDown" ? 1 : -1;
      setActive(Math.max(0, Math.min(shown.length - 1, active + step)));
    } else if (e.key === "Enter" && !list.hidden) {
      // Choose the highlighted font instead of submitting the form
      e.preventDefault();
      if (active >= 0) select(shown[active].file);
      else commit();
    } else if (e.key === "Enter") {
      // Let the form submit with the typed font
      commit();
    } else if (e.key === "Escape" && !list.hidden) {
      // Close the list without closing the window
      e.preventDefault();
      e.stopPropagation();
      close();
      sync();
    }
  });
  window.addEventListener("resize", close);

  return {
    sync,
    close,
    setFonts(list) {
      fonts = [embedded, ...list];
      sync();
    },
  };
}
