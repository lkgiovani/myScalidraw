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

export interface ExcalidrawData {
  elements: ExcalidrawElement[];
  appState: ExcalidrawAppState;
}
