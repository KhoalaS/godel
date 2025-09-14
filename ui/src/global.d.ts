/**
 * Augment Window interface, adding showOpenFilePicker function.
 */
declare global {
  interface Window {
    showOpenFilePicker: (arg: {
      types: { description: string; accept: Record<string, string[]> }[]
    }) => Promise<FileSystemFileHandle[]>
  }
}

export {}
