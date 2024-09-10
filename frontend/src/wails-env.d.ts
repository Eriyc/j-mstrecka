import * as Runtime from "../wailsjs/runtime/runtime";

declare global {
  interface Window {
    runtime: typeof Runtime;
  }
}
