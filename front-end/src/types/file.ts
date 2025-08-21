export interface ExcalidrawFile {
  id: string;
  name: string;
  data: ExcalidrawData; // Excalidraw scene data
  lastModified: number;
  parentId?: string; // For folder structure
}

export interface ExcalidrawData {
  elements: ExcalidrawElement[];
  appState: ExcalidrawAppState;
}

export interface ExcalidrawElement {
  id: string;
  type: string;
  x: number;
  y: number;
  width: number;
  height: number;
  angle: number;
  strokeColor: string;
  backgroundColor: string;
  fillStyle: string;
  strokeWidth: number;
  strokeStyle: string;
  roughness: number;
  opacity: number;
  [key: string]: unknown;
}

export interface ExcalidrawAppState {
  viewBackgroundColor: string;
  currentItemStrokeColor: string;
  currentItemBackgroundColor: string;
  currentItemFillStyle: string;
  currentItemStrokeWidth: number;
  currentItemStrokeStyle: string;
  currentItemRoughness: number;
  currentItemOpacity: number;
  currentItemFontFamily: number;
  currentItemFontSize: number;
  currentItemTextAlign: string;
  currentItemStartArrowhead: string | null;
  currentItemEndArrowhead: string | null;
  scrollX: number;
  scrollY: number;
  zoom: {
    value: number;
  };
  currentItemRoundness: string;
  gridSize: number | null;
  colorPalette: Record<string, string>;
  [key: string]: unknown;
}

export interface FileFolder {
  id: string;
  name: string;
  isFolder: true;
  parentId?: string;
  children: (ExcalidrawFile | FileFolder)[];
  isExpanded?: boolean;
}

export type FileSystemItem = ExcalidrawFile | FileFolder;

export interface FileContextMenuProps {
  x: number;
  y: number;
  item: FileSystemItem;
  onRename: (id: string, newName: string) => void;
  onDelete: (id: string) => void;
  onClose: () => void;
}
