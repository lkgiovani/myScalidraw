import { Excalidraw } from "@excalidraw/excalidraw";

import "@excalidraw/excalidraw/index.css";
import { useFileStore } from "@/stores/fileStore";
import { useFile, useSaveFile } from "@/hooks/useFileApi";
import {
  useState,
  useCallback,
  useEffect,
  useRef,
  Component,
  ReactNode,
} from "react";
import { toast } from "sonner";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type ExcalidrawData = any;

interface FileData {
  content?: string;
  name?: string;
  data?: Record<string, unknown>;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

class ErrorBoundary extends Component<
  { children: ReactNode },
  ErrorBoundaryState
> {
  constructor(props: { children: ReactNode }) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: unknown) {
    console.error("Excalidraw Error:", error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="flex items-center justify-center h-full bg-background">
          <div className="text-center space-y-4">
            <div className="w-16 h-16 mx-auto bg-red-100 rounded-full flex items-center justify-center">
              <svg
                className="w-8 h-8 text-red-500"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
                />
              </svg>
            </div>
            <div>
              <h3 className="text-lg font-medium text-foreground">
                Erro ao renderizar
              </h3>
              <p className="text-sm text-muted-foreground">
                Ocorreu um erro ao carregar o editor. Tente recarregar a página.
              </p>
              <button
                onClick={() => window.location.reload()}
                className="mt-2 px-4 py-2 bg-primary text-primary-foreground rounded-md text-sm hover:bg-primary/90"
              >
                Recarregar
              </button>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export function ExcalidrawEditor() {
  const selectedFileId = useFileStore((state) => state.selectedFileId);
  const selectedFile = useFileStore((state) =>
    selectedFileId ? state.files[selectedFileId] : null
  );

  const { data: fileData, isLoading: isLoadingFile } = useFile(selectedFileId);
  const saveFileMutation = useSaveFile();

  const [excalidrawData, setExcalidrawData] = useState<ExcalidrawData | null>(
    null
  );

  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const autoSaveTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    // Resetar dados quando mudar de arquivo ou não tiver dados
    if (!fileData || !selectedFileId) {
      setExcalidrawData(null);
      return;
    }

    const typedFileData = fileData as FileData;

    let parsedContent = null;

    if (typedFileData?.content) {
      try {
        parsedContent = JSON.parse(typedFileData.content);
      } catch (error) {
        console.error("Erro ao parsear conteúdo string:", error);
      }
    } else if (typedFileData?.data) {
      parsedContent = typedFileData.data;
    }

    if (parsedContent) {
      const elements = Array.isArray(parsedContent.elements)
        ? parsedContent.elements
        : [];
      const appState = parsedContent.appState || {};

      const sanitizedAppState = {
        viewBackgroundColor: appState.viewBackgroundColor || "#ffffff",
        gridSize: appState.gridSize || null,
        gridStep: appState.gridStep || 5,
        gridModeEnabled: appState.gridModeEnabled || false,
        zoom:
          appState.zoom && typeof appState.zoom === "object"
            ? appState.zoom
            : { value: 1 },
        scrollX: typeof appState.scrollX === "number" ? appState.scrollX : 0,
        scrollY: typeof appState.scrollY === "number" ? appState.scrollY : 0,
        ...appState,

        collaborators: appState.collaborators
          ? new Map(Object.entries(appState.collaborators))
          : new Map(),
      };

      const newExcalidrawData = {
        elements,
        appState: sanitizedAppState,
      };

      setExcalidrawData(newExcalidrawData);
      setHasUnsavedChanges(false);
    }
  }, [fileData, selectedFileId]);

  const autoSaveRef =
    useRef<
      (elements: ExcalidrawData, appState: ExcalidrawData) => Promise<void>
    >();

  autoSaveRef.current = async (
    elements: ExcalidrawData,
    appState: ExcalidrawData
  ) => {
    if (!selectedFileId) return;

    setIsSaving(true);

    const jsonData = {
      type: "excalidraw",
      version: 2,
      source: "https://excalidraw.com",
      elements,
      appState,
      files: {},
    };

    const contentToSave = JSON.stringify(jsonData, null, 2);

    try {
      await saveFileMutation.mutateAsync({
        id: selectedFileId,
        data: { content: contentToSave },
      });
      setHasUnsavedChanges(false);
    } catch (error) {
      console.error("Erro no auto-save:", error);
      toast.error("Erro ao salvar automaticamente");
    } finally {
      setIsSaving(false);
    }
  };

  const handleChange = useCallback(
    (elements: ExcalidrawData, appState: ExcalidrawData) => {
      setExcalidrawData({ elements, appState });
      setHasUnsavedChanges(true);

      if (autoSaveTimeoutRef.current) {
        clearTimeout(autoSaveTimeoutRef.current);
      }

      if (selectedFileId && autoSaveRef.current) {
        autoSaveTimeoutRef.current = setTimeout(() => {
          autoSaveRef.current?.(elements, appState);
        }, 2000);
      }
    },
    [selectedFileId]
  );

  const handleSave = useCallback(() => {
    if (!selectedFileId || !hasUnsavedChanges || !excalidrawData) return;

    const jsonData = {
      type: "excalidraw",
      version: 2,
      source: "https://excalidraw.com",
      elements: excalidrawData.elements,
      appState: excalidrawData.appState,
      files: {},
    };

    const contentToSave = JSON.stringify(jsonData, null, 2);

    saveFileMutation.mutate(
      { id: selectedFileId, data: { content: contentToSave } },
      {
        onSuccess: () => {
          setHasUnsavedChanges(false);
        },
      }
    );
  }, [selectedFileId, excalidrawData, hasUnsavedChanges, saveFileMutation]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "s") {
        e.preventDefault();
        handleSave();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => {
      window.removeEventListener("keydown", handleKeyDown);

      if (autoSaveTimeoutRef.current) {
        clearTimeout(autoSaveTimeoutRef.current);
      }
    };
  }, [handleSave]);

  if (isLoadingFile) {
    return (
      <div className="flex-1 flex items-center justify-center bg-explorer-main">
        <div className="text-center space-y-4">
          <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto"></div>
          <p className="text-sm text-muted-foreground">Carregando arquivo...</p>
        </div>
      </div>
    );
  }

  if (!selectedFile) {
    return (
      <div className="flex-1 flex items-center justify-center bg-explorer-main">
        <div className="text-center space-y-4">
          <div className="w-24 h-24 mx-auto bg-muted rounded-lg flex items-center justify-center">
            <svg
              className="w-12 h-12 text-muted-foreground"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
              />
            </svg>
          </div>
          <div>
            <h3 className="text-lg font-medium text-foreground">
              Selecione um arquivo
            </h3>
            <p className="text-sm text-muted-foreground">
              Escolha um arquivo .excalidraw no explorador para começar a editar
            </p>
          </div>
        </div>
      </div>
    );
  }

  if (selectedFile.type === "folder") {
    return (
      <div className="flex-1 flex items-center justify-center bg-explorer-main">
        <div className="text-center space-y-4">
          <div className="w-24 h-24 mx-auto bg-muted rounded-lg flex items-center justify-center">
            <svg
              className="w-12 h-12 text-muted-foreground"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25H11.69Z"
              />
            </svg>
          </div>
          <div>
            <h3 className="text-lg font-medium text-foreground">
              Folder selected
            </h3>
            <p className="text-sm text-muted-foreground">
              Please select a file to edit
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 bg-explorer-main relative">
      {/* Botão de Salvar no canto superior direito */}
      {selectedFileId && (
        <div className="absolute top-4 right-4 z-20 flex items-center gap-2">
          <button
            onClick={handleSave}
            disabled={
              !hasUnsavedChanges || saveFileMutation.isPending || isSaving
            }
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white rounded-md text-sm font-medium transition-colors flex items-center gap-2"
          >
            {saveFileMutation.isPending || isSaving ? (
              <>
                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                Salvando...
              </>
            ) : (
              <>
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3-3m0 0l-3 3m3-3v12"
                  />
                </svg>
                Salvar
              </>
            )}
          </button>
        </div>
      )}

      {/* Status indicators */}
      {isSaving && (
        <div className="absolute top-16 right-4 z-10 bg-blue-500 text-white px-3 py-1 rounded-md text-sm font-medium">
          Auto-salvando...
        </div>
      )}

      {hasUnsavedChanges && !isSaving && (
        <div className="absolute top-16 right-4 z-10 bg-yellow-500 text-yellow-900 px-3 py-1 rounded-md text-sm font-medium">
          Alterações não salvas
        </div>
      )}

      {!hasUnsavedChanges &&
        !isSaving &&
        selectedFileId &&
        !saveFileMutation.isPending && (
          <div className="absolute top-16 right-4 z-10 bg-green-500 text-white px-3 py-1 rounded-md text-sm font-medium">
            ✓ Salvo
          </div>
        )}

      <div className="h-full w-full">
        <div style={{ height: "100%", width: "100%" }}>
          {/* Só renderizar Excalidraw quando tiver dados válidos do backend */}
          {!isLoadingFile &&
          selectedFileId &&
          fileData &&
          excalidrawData &&
          excalidrawData.elements &&
          Array.isArray(excalidrawData.elements) ? (
            <ErrorBoundary>
              <Excalidraw
                key={selectedFileId || "default"}
                initialData={excalidrawData}
                onChange={handleChange}
                theme="light"
                UIOptions={{
                  canvasActions: {
                    loadScene: false,
                    toggleTheme: true,
                  },
                }}
                viewModeEnabled={false}
                zenModeEnabled={false}
                gridModeEnabled={false}
              />
            </ErrorBoundary>
          ) : (
            <div className="flex items-center justify-center h-full text-gray-500">
              {isLoadingFile ? (
                <div className="flex items-center gap-2">
                  <div className="w-6 h-6 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
                  Carregando arquivo do servidor...
                </div>
              ) : selectedFileId && !fileData ? (
                <div className="flex items-center gap-2">
                  <div className="w-6 h-6 border-2 border-orange-500 border-t-transparent rounded-full animate-spin"></div>
                  Aguardando dados do backend...
                </div>
              ) : selectedFileId &&
                fileData &&
                (!excalidrawData ||
                  !excalidrawData.elements ||
                  !Array.isArray(excalidrawData.elements)) ? (
                <div className="flex items-center gap-2">
                  <div className="w-6 h-6 border-2 border-purple-500 border-t-transparent rounded-full animate-spin"></div>
                  Processando conteúdo do arquivo...
                </div>
              ) : (
                "Selecione um arquivo para editar"
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
