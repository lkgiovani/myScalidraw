import { useState, useCallback, useEffect, useRef } from "react";
import { useFileStore } from "@/stores/fileStore";
import { useFile, useSaveFile } from "@/hooks/useFileApi";
import { toast } from "sonner";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type ExcalidrawData = any;

interface FileData {
  content?: string;
  name?: string;
  data?: Record<string, unknown>;
}

export function useExcalidrawEditor() {
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
  const [currentFileId, setCurrentFileId] = useState<string | null>(null);
  const autoSaveTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Parse file data and update excalidraw state
  useEffect(() => {
    if (autoSaveTimeoutRef.current) {
      clearTimeout(autoSaveTimeoutRef.current);
      autoSaveTimeoutRef.current = null;
    }

    if (selectedFileId !== currentFileId) {
      setCurrentFileId(selectedFileId);
      setExcalidrawData(null);
      setHasUnsavedChanges(false);
      setIsSaving(false);
    }

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

      if (selectedFileId === currentFileId) {
        setHasUnsavedChanges(false);
      }
    }
  }, [fileData, selectedFileId, currentFileId]);

  // Auto-save function
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

  // Handle changes in excalidraw
  const handleChange = useCallback(
    (elements: ExcalidrawData, appState: ExcalidrawData) => {
      if (!selectedFileId || selectedFileId !== currentFileId) {
        return;
      }

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
    [selectedFileId, currentFileId]
  );

  // Manual save function
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

  // Keyboard shortcuts
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
        autoSaveTimeoutRef.current = null;
      }
    };
  }, [handleSave]);

  // Cleanup on file change
  useEffect(() => {
    return () => {
      if (autoSaveTimeoutRef.current) {
        clearTimeout(autoSaveTimeoutRef.current);
        autoSaveTimeoutRef.current = null;
      }
    };
  }, [selectedFileId]);

  return {
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
  };
}
