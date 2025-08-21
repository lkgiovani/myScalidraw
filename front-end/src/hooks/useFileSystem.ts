import { useState, useCallback, useEffect } from "react";
import { ExcalidrawFile, FileFolder, FileSystemItem } from "@/types/file";

// Função para carregar arquivo Excalidraw do public
const loadExcalidrawFile = async (filename: string): Promise<any> => {
  try {
    const response = await fetch(`/${filename}`);
    if (!response.ok) {
      throw new Error(`Failed to load ${filename}`);
    }
    return await response.json();
  } catch (error) {
    console.error("Error loading Excalidraw file:", error);
    return null;
  }
};

export const useFileSystem = () => {
  const [files, setFiles] = useState<FileSystemItem[]>([
    {
      id: "rascunhos",
      name: "Rascunhos",
      isFolder: true,
      children: [
        {
          id: "exemplo-salve",
          name: "Untitled-2025-06-30-1107",
          data: { elements: [], appState: { viewBackgroundColor: "#ffffff" } },
          lastModified: Date.now() - 3600000,
          parentId: "rascunhos",
        } as ExcalidrawFile,
      ],
      isExpanded: true,
    } as FileFolder,
  ]);

  const [activeFileId, setActiveFileId] = useState<string>("");

  // Carregar arquivo do public na inicialização
  useEffect(() => {
    const loadInitialFile = async () => {
      const excalidrawData = await loadExcalidrawFile(
        "Untitled-2025-06-30-1107.excalidraw"
      );
      if (excalidrawData) {
        setFiles((prev) =>
          prev.map((item) => {
            if ("isFolder" in item && item.id === "rascunhos") {
              return {
                ...item,
                children: item.children.map((child) =>
                  child.id === "exemplo-salve"
                    ? {
                        ...child,
                        data: {
                          elements: excalidrawData.elements,
                          appState: excalidrawData.appState,
                        },
                      }
                    : child
                ),
              };
            }
            return item;
          })
        );
      }
    };

    loadInitialFile();
  }, []);

  const createFile = useCallback((name: string, parentId?: string) => {
    const newFile: ExcalidrawFile = {
      id: `file-${Date.now()}`,
      name,
      data: { elements: [], appState: { viewBackgroundColor: "#ffffff" } },
      lastModified: Date.now(),
      parentId,
    };

    setFiles((prev) => [...prev, newFile]);
    setActiveFileId(newFile.id);
    return newFile.id;
  }, []);

  const createFolder = useCallback((name: string, parentId?: string) => {
    const newFolder: FileFolder = {
      id: `folder-${Date.now()}`,
      name,
      isFolder: true,
      parentId,
      children: [],
      isExpanded: true,
    };

    setFiles((prev) => [...prev, newFolder]);
    return newFolder.id;
  }, []);

  const updateFile = useCallback((id: string, data: any) => {
    setFiles((prev) =>
      prev.map((file) =>
        file.id === id && !("isFolder" in file)
          ? ({ ...file, data, lastModified: Date.now() } as ExcalidrawFile)
          : file
      )
    );
  }, []);

  const renameItem = useCallback((id: string, newName: string) => {
    setFiles((prev) =>
      prev.map((item) => (item.id === id ? { ...item, name: newName } : item))
    );
  }, []);

  const deleteItem = useCallback(
    (id: string) => {
      setFiles((prev) => prev.filter((item) => item.id !== id));
      if (activeFileId === id) {
        // Buscar o primeiro arquivo disponível
        const findFirstFile = (items: FileSystemItem[]): string | null => {
          for (const item of items) {
            if (!("isFolder" in item)) {
              return item.id;
            } else if (item.children.length > 0) {
              const childFile = findFirstFile(item.children);
              if (childFile) return childFile;
            }
          }
          return null;
        };

        const remainingFiles = files.filter((f) => f.id !== id);
        const firstFileId = findFirstFile(remainingFiles);

        if (firstFileId) {
          setActiveFileId(firstFileId);
        } else {
          setActiveFileId("");
        }
      }
    },
    [activeFileId, files]
  );

  const toggleFolder = useCallback((id: string) => {
    setFiles((prev) =>
      prev.map((item) =>
        item.id === id && "isFolder" in item
          ? ({ ...item, isExpanded: !item.isExpanded } as FileFolder)
          : item
      )
    );
  }, []);

  const getActiveFile = useCallback(() => {
    const findFileById = (
      items: FileSystemItem[],
      id: string
    ): ExcalidrawFile | undefined => {
      for (const item of items) {
        if (item.id === id && !("isFolder" in item)) {
          return item as ExcalidrawFile;
        } else if ("isFolder" in item && item.children.length > 0) {
          const found = findFileById(item.children, id);
          if (found) return found;
        }
      }
      return undefined;
    };

    return findFileById(files, activeFileId);
  }, [files, activeFileId]);

  return {
    files,
    activeFileId,
    setActiveFileId,
    createFile,
    createFolder,
    updateFile,
    renameItem,
    deleteItem,
    toggleFolder,
    getActiveFile,
  };
};
