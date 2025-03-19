export {}

declare global {
  interface Window {
    runtime: {
      EventsOn: (event: string, callback: (data: string) => void) => void
    }
  }
}
