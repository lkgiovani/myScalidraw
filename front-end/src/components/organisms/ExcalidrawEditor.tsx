import { Excalidraw } from "@excalidraw/excalidraw";
import "@excalidraw/excalidraw/index.css";

import { useExcalidrawEditor } from "@/hooks/useExcalidrawEditor";
import { ErrorBoundary } from "@/components/atoms/ErrorBoundary";
import { LoadingState } from "@/components/molecules/LoadingState";
import { EmptyState } from "@/components/molecules/EmptyState";
import { SaveButton } from "@/components/molecules/SaveButton";
import { StatusIndicators } from "@/components/molecules/StatusIndicators";
import { FileText, Folder } from "lucide-react";
import { useTheme } from "@/contexts/ThemeContext";

export function ExcalidrawEditor() {
  const { theme } = useTheme();
  const {
    selectedFileId,
    selectedFile,
    fileData,
    isLoadingFile,
    excalidrawData,
    hasUnsavedChanges,
    isSaving,
    saveFileMutation,
    handleChange,
    handleSave,
  } = useExcalidrawEditor();

  if (isLoadingFile) {
    return <LoadingState message="Carregando arquivo..." />;
  }

  if (!selectedFile) {
    return (
      <EmptyState
        icon={
          <FileText
            className="w-12 h-12 text-muted-foreground"
            strokeWidth={1.5}
          />
        }
        title="Selecione um arquivo"
        description="Escolha um arquivo .excalidraw no explorador para começar a editar"
      />
    );
  }

  if (selectedFile.type === "folder") {
    return (
      <EmptyState
        icon={
          <Folder
            className="w-12 h-12 text-muted-foreground"
            strokeWidth={1.5}
          />
        }
        title="Pasta selecionada"
        description="Selecione um arquivo para editar"
      />
    );
  }

  return (
    <div className="flex-1 bg-explorer-main relative">
      {/* Save button */}
      {selectedFileId && (
        <div className="absolute top-4 right-4 z-20">
          <SaveButton
            onSave={handleSave}
            hasUnsavedChanges={hasUnsavedChanges}
            isSaving={isSaving || saveFileMutation.isPending}
          />
        </div>
      )}

      {/* Status indicators */}
      <StatusIndicators
        isSaving={isSaving}
        hasUnsavedChanges={hasUnsavedChanges}
        selectedFileId={selectedFileId}
        isPending={saveFileMutation.isPending}
      />

      {/* Excalidraw canvas */}
      <div className="h-full w-full">
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
              theme={theme === "system" ? undefined : theme}
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
          <div className="flex items-center justify-center h-full">
            {isLoadingFile ? (
              <LoadingState
                message="Carregando arquivo do servidor..."
                spinnerColor="blue"
              />
            ) : selectedFileId && !fileData ? (
              <LoadingState
                message="Aguardando dados do backend..."
                spinnerColor="orange"
              />
            ) : selectedFileId &&
              fileData &&
              (!excalidrawData ||
                !excalidrawData.elements ||
                !Array.isArray(excalidrawData.elements)) ? (
              <LoadingState
                message="Processando conteúdo do arquivo..."
                spinnerColor="purple"
              />
            ) : (
              <p className="text-muted-foreground">
                Selecione um arquivo para editar
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
