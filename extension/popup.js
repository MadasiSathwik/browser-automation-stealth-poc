document.getElementById("start").addEventListener("click", () => {
  console.log("[EXTENSION] Demo Mode Enabled");
  console.log("[EXTENSION] Safe Mode Enabled");

  chrome.runtime.sendMessage({
    type: "START_DEMO",
    demo: true,
    safe: true
  });
});
