import { create } from "zustand";
import { fileApi, type FileMetadata } from "@/lib/api";

const getDefaultExcalidrawContent = async (): Promise<string> => {
  try {
    const response = await fetch("/defund.excalidraw");
    const content = await response.text();
    return content;
  } catch (error) {
    console.error("Erro ao carregar template padrão:", error);
    return JSON.stringify({
      type: "excalidraw",
      version: 2,
      source: "https://excalidraw.com",
      elements: [],
      appState: {
        gridSize: 20,
        gridStep: 5,
        gridModeEnabled: false,
        viewBackgroundColor: "#ffffff",
        lockedMultiSelections: {},
      },
      files: {},
    });
  }
};

export interface FileNode {
  id: string;
  name: string;
  type: "file" | "folder";
  path: string;
  parentId?: string;
  children?: string[];
  size?: number;
  modified?: Date;
  icon?: string;
  isOpen?: boolean;
  content?: string;
}

interface FileStore {
  files: Record<string, FileNode>;
  rootFolders: string[];

  selectedFileId: string | null;
  expandedFolders: Set<string>;
  searchQuery: string;
  isLoading: boolean;

  setFiles: (files: Record<string, FileNode>) => void;
  loadFilesFromApi: () => Promise<void>;
  selectFile: (fileId: string | null) => void;
  toggleFolder: (folderId: string) => void;
  setSearchQuery: (query: string) => void;
  setLoading: (loading: boolean) => void;

  createFileApi: (name: string, parentId?: string) => Promise<void>;
  createFolderApi: (name: string, parentId?: string) => Promise<void>;
  saveFileApi: (id: string, content: string) => Promise<void>;
  deleteFileApi: (id: string) => Promise<void>;

  createFile: (name: string, parentId?: string) => string;
  createFolder: (name: string, parentId?: string) => string;
  deleteItem: (itemId: string) => void;
  renameItem: (itemId: string, newName: string) => void;
}

const convertApiFileToNode = (apiFile: FileMetadata): FileNode => {
  let content: string | undefined;

  if (apiFile.content) {
    content = apiFile.content;
  } else if (apiFile.data) {
    content = JSON.stringify(apiFile.data);
  }

  const fileName = apiFile.name || apiFile.fileName;

  const node: FileNode = {
    id: apiFile.id,
    name: fileName,
    type: apiFile.isFolder ? "folder" : "file",
    path: apiFile.path || `/${fileName}`,
    size: apiFile.size,
    modified: apiFile.modified
      ? new Date(apiFile.modified)
      : apiFile.lastModified
      ? new Date(apiFile.lastModified)
      : undefined,
    content,
    parentId: apiFile.parentId,
  };

  return node;
};

const buildFileTree = (
  files: FileMetadata[]
): { filesMap: Record<string, FileNode>; rootFiles: string[] } => {
  const filesMap: Record<string, FileNode> = {};
  const rootFiles: string[] = [];

  files.forEach((file) => {
    const node = convertApiFileToNode(file);
    filesMap[node.id] = node;
  });

  Object.values(filesMap).forEach((node) => {
    if (!node.parentId || node.parentId === "") {
      rootFiles.push(node.id);
    } else {
      const parent = filesMap[node.parentId];
      if (parent) {
        if (!parent.children) {
          parent.children = [];
        }
        parent.children.push(node.id);
      } else {
        rootFiles.push(node.id);
      }
    }
  });

  return { filesMap, rootFiles };
};

export const useFileStore = create<FileStore>((set, get) => ({
  files: {},
  rootFolders: [],
  selectedFileId: null,
  expandedFolders: new Set(),
  searchQuery: "",
  isLoading: false,

  setFiles: (files) => set({ files }),

  setLoading: (loading) => set({ isLoading: loading }),

  loadFilesFromApi: async () => {
    try {
      set({ isLoading: true });

      const apiFiles = await fileApi.getFiles();

      const { filesMap, rootFiles } = buildFileTree(apiFiles);

      set({
        files: filesMap,
        rootFolders: rootFiles,
        isLoading: false,
      });
    } catch (error) {
      console.error("Erro ao carregar arquivos:", error);
      set({ isLoading: false });
    }
  },

  selectFile: (fileId) => set({ selectedFileId: fileId }),

  toggleFolder: (folderId) =>
    set((state) => {
      const newExpanded = new Set(state.expandedFolders);
      if (newExpanded.has(folderId)) {
        newExpanded.delete(folderId);
      } else {
        newExpanded.add(folderId);
      }
      return { expandedFolders: newExpanded };
    }),

  setSearchQuery: (query) => set({ searchQuery: query }),

  createFileApi: async (name, parentId) => {
    try {
      const defaultContent = await getDefaultExcalidrawContent();
      const newFile = await fileApi.createFile({
        name,
        parentId,
        content: defaultContent,
        isFolder: false,
      });

      const node = convertApiFileToNode(newFile);
      node.parentId = parentId;

      set((state) => {
        const updatedFiles = { ...state.files, [node.id]: node };

        if (parentId && state.files[parentId]) {
          const parent = state.files[parentId];
          updatedFiles[parentId] = {
            ...parent,
            children: [...(parent.children || []), node.id],
          };
          return { files: updatedFiles };
        }

        return {
          files: updatedFiles,
          rootFolders: [...state.rootFolders, node.id],
        };
      });
    } catch (error) {
      console.error("Erro ao criar arquivo:", error);
      throw error;
    }
  },

  createFolderApi: async (name, parentId) => {
    try {
      const newFolder = await fileApi.createFile({
        name,
        parentId,
        isFolder: true,
      });

      const node = convertApiFileToNode(newFolder);
      node.parentId = parentId;
      node.children = [];

      set((state) => {
        const updatedFiles = { ...state.files, [node.id]: node };

        if (parentId && state.files[parentId]) {
          const parent = state.files[parentId];
          updatedFiles[parentId] = {
            ...parent,
            children: [...(parent.children || []), node.id],
          };
          return { files: updatedFiles };
        }

        return {
          files: updatedFiles,
          rootFolders: [...state.rootFolders, node.id],
        };
      });
    } catch (error) {
      console.error("Erro ao criar pasta:", error);
      throw error;
    }
  },

  saveFileApi: async (id, content) => {
    try {
      await fileApi.saveFile(id, { content });

      set((state) => ({
        files: {
          ...state.files,
          [id]: { ...state.files[id], content },
        },
      }));
    } catch (error) {
      console.error("Erro ao salvar arquivo:", error);
      throw error;
    }
  },

  deleteFileApi: async (id) => {
    try {
      await fileApi.deleteFile(id);

      set((state) => {
        const newFiles = { ...state.files };
        delete newFiles[id];

        return {
          files: newFiles,
          rootFolders: state.rootFolders.filter((rootId) => rootId !== id),
          selectedFileId:
            state.selectedFileId === id ? null : state.selectedFileId,
        };
      });
    } catch (error) {
      console.error("Erro ao excluir arquivo:", error);
      throw error;
    }
  },

  createFile: (name, parentId) => {
    const id = `file-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    const now = new Date();

    set((state) => {
      const newFile: FileNode = {
        id,
        name: name.endsWith(".excalidraw") ? name : `${name}.excalidraw`,
        type: "file",
        path: parentId ? `${state.files[parentId]?.path}/${name}` : `/${name}`,
        parentId,
        size: 0,
        modified: now,
      };

      const updatedFiles = { ...state.files, [id]: newFile };

      if (parentId && state.files[parentId]) {
        const parent = state.files[parentId];
        updatedFiles[parentId] = {
          ...parent,
          children: [...(parent.children || []), id],
        };
      } else {
        return {
          files: updatedFiles,
          rootFolders: [...state.rootFolders, id],
        };
      }

      return { files: updatedFiles };
    });

    return id;
  },

  createFolder: (name, parentId) => {
    const id = `folder-${Date.now()}-${Math.random()
      .toString(36)
      .substr(2, 9)}`;

    set((state) => {
      const newFolder: FileNode = {
        id,
        name,
        type: "folder",
        path: parentId ? `${state.files[parentId]?.path}/${name}` : `/${name}`,
        parentId,
        children: [],
        isOpen: false,
      };

      const updatedFiles = { ...state.files, [id]: newFolder };

      if (parentId && state.files[parentId]) {
        const parent = state.files[parentId];
        updatedFiles[parentId] = {
          ...parent,
          children: [...(parent.children || []), id],
        };
      } else {
        return {
          files: updatedFiles,
          rootFolders: [...state.rootFolders, id],
        };
      }

      return { files: updatedFiles };
    });

    return id;
  },

  deleteItem: (itemId) => {
    set((state) => {
      const item = state.files[itemId];
      if (!item) return state;

      const updatedFiles = { ...state.files };
      const itemsToDelete = [itemId];

      if (item.type === "folder" && item.children) {
        const collectChildren = (children: string[]) => {
          children.forEach((childId) => {
            itemsToDelete.push(childId);
            const child = updatedFiles[childId];
            if (child?.type === "folder" && child.children) {
              collectChildren(child.children);
            }
          });
        };
        collectChildren(item.children);
      }

      itemsToDelete.forEach((id) => {
        delete updatedFiles[id];
      });

      let updatedRootFolders = state.rootFolders;

      if (item.parentId && updatedFiles[item.parentId]) {
        const parent = updatedFiles[item.parentId];
        updatedFiles[item.parentId] = {
          ...parent,
          children:
            parent.children?.filter((childId) => childId !== itemId) || [],
        };
      } else {
        updatedRootFolders = state.rootFolders.filter((id) => id !== itemId);
      }

      return {
        files: updatedFiles,
        rootFolders: updatedRootFolders,
        selectedFileId:
          state.selectedFileId === itemId ? null : state.selectedFileId,
      };
    });
  },

  renameItem: (itemId, newName) => {
    set((state) => {
      const item = state.files[itemId];
      if (!item) return state;

      const updatedFiles = { ...state.files };
      const finalName =
        item.type === "file" && !newName.endsWith(".excalidraw")
          ? `${newName}.excalidraw`
          : newName;

      updatedFiles[itemId] = {
        ...item,
        name: finalName,
        path: item.parentId
          ? `${state.files[item.parentId]?.path}/${finalName}`
          : `/${finalName}`,
      };

      return { files: updatedFiles };
    });
  },
}));
