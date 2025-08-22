import { useState, useEffect, useCallback } from "react";
import {
  FileEntity,
  FileSystemItemEntity,
} from "../../domain/entities/FileEntity";
import { ExcalidrawData } from "../../domain/entities/ExcalidrawEntity";
import { FileController } from "../controllers/FileController";

export const useFileSystemPresenter = (fileController: FileController) => {
  const [files, setFiles] = useState<FileSystemItemEntity[]>([]);
  const [activeFileId, setActiveFileId] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(true);

  // Load initial files
  useEffect(() => {
    const loadFiles = async () => {
      try {
        setIsLoading(true);
        const filesData = await fileController.getFiles();
        setFiles(filesData);

        // Load initial file from public directory
        const excalidrawData = await fileController.loadExcalidrawFile(
          "Untitled-2025-06-30-1107.excalidraw"
        );
        if (excalidrawData) {
          // Find the example file and update it with the loaded data
          const exampleFile = filesData
            .flatMap((item) => ("isFolder" in item ? item.children : [item]))
            .find((item) => item.id === "exemplo-salve");

          if (exampleFile) {
            await fileController.updateFile(exampleFile.id, excalidrawData);
            // Refresh files after update
            const updatedFiles = await fileController.getFiles();
            setFiles(updatedFiles);
          }
        }
      } catch (error) {
        console.error("Error loading files:", error);
      } finally {
        setIsLoading(false);
      }
    };

    loadFiles();
  }, [fileController]);

  const createFile = useCallback(
    async (name: string, parentId?: string) => {
      try {
        const fileId = await fileController.createFile(name, parentId);
        const updatedFiles = await fileController.getFiles();
        setFiles(updatedFiles);
        setActiveFileId(fileId);
        return fileId;
      } catch (error) {
        console.error("Error creating file:", error);
        throw error;
      }
    },
    [fileController]
  );

  const createFolder = useCallback(
    async (name: string, parentId?: string) => {
      try {
        const folderId = await fileController.createFolder(name, parentId);
        const updatedFiles = await fileController.getFiles();
        setFiles(updatedFiles);
        return folderId;
      } catch (error) {
        console.error("Error creating folder:", error);
        throw error;
      }
    },
    [fileController]
  );

  const updateFile = useCallback(
    async (id: string, data: ExcalidrawData) => {
      try {
        await fileController.updateFile(id, data);
        const updatedFiles = await fileController.getFiles();
        setFiles(updatedFiles);
      } catch (error) {
        console.error("Error updating file:", error);
      }
    },
    [fileController]
  );

  const renameItem = useCallback(
    async (id: string, newName: string) => {
      try {
        await fileController.renameItem(id, newName);
        const updatedFiles = await fileController.getFiles();
        setFiles(updatedFiles);
      } catch (error) {
        console.error("Error renaming item:", error);
      }
    },
    [fileController]
  );

  const deleteItem = useCallback(
    async (id: string) => {
      try {
        await fileController.deleteItem(id);
        const updatedFiles = await fileController.getFiles();
        setFiles(updatedFiles);

        if (activeFileId === id) {
          // Find another file to select
          const findFirstFile = (
            items: FileSystemItemEntity[]
          ): string | null => {
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

          const firstFileId = findFirstFile(updatedFiles);
          setActiveFileId(firstFileId || "");
        }
      } catch (error) {
        console.error("Error deleting item:", error);
      }
    },
    [fileController, activeFileId]
  );

  const toggleFolder = useCallback(
    async (id: string) => {
      try {
        await fileController.toggleFolder(id);
        const updatedFiles = await fileController.getFiles();
        setFiles(updatedFiles);
      } catch (error) {
        console.error("Error toggling folder:", error);
      }
    },
    [fileController]
  );

  const getActiveFile = useCallback(async (): Promise<
    FileEntity | undefined
  > => {
    if (!activeFileId) return undefined;
    try {
      return await fileController.getFileById(activeFileId);
    } catch (error) {
      console.error("Error getting active file:", error);
      return undefined;
    }
  }, [fileController, activeFileId]);

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
    isLoading,
  };
};
