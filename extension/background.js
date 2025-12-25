chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === "START_DEMO") {
    console.log("[BACKGROUND] Starting automation demo...");
    console.log("[ACTION] Initializing stealth profile");
    console.log("[ACTION] Simulating mouse movement (Bezier)");
    console.log("[ACTION] Typing with human delay");
    console.log("[ACTION] Stealth techniques applied: 8");
  }
});
