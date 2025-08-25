import ky from "ky";

const API_BASE_URL = "http://localhost:8181/api";

export const api = ky.create({
  prefixUrl: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 30000,
});

export interface FileMetadata {
  id: string;
  name: string;
  fileName?: string; // Backend pode retornar fileName ao invés de name
  path?: string;
  size?: number;
  modified?: string;
  lastModified?: number;
  content?: string;
  data?: Record<string, unknown>;
  isFolder?: boolean;
  parentId?: string;
  children?: FileMetadata[]; // Para estrutura hierárquica do backend
}

export interface CreateFileRequest {
  name: string;
  parentId?: string;
  content?: string;
  isFolder?: boolean;
}

export interface SaveFileRequest {
  content: string;
}

export interface UploadFileRequest {
  file: File;
  parentPath?: string;
}

export const fileApi = {
  getFiles: async (): Promise<FileMetadata[]> => {
    return api.get("files").json<FileMetadata[]>();
  },

  getFileById: async (id: string): Promise<FileMetadata> => {
    return api.get(`files/${id}`).json<FileMetadata>();
  },

  createFile: async (data: CreateFileRequest): Promise<FileMetadata> => {
    return api.post("files", { json: data }).json<FileMetadata>();
  },

  uploadFile: async (data: UploadFileRequest): Promise<FileMetadata> => {
    const formData = new FormData();
    formData.append("file", data.file);
    if (data.parentPath) {
      formData.append("parentPath", data.parentPath);
    }

    return api.post("files/upload", { body: formData }).json<FileMetadata>();
  },

  saveFile: async (
    id: string,
    data: SaveFileRequest
  ): Promise<{ message: string }> => {
    return api
      .put(`files/${id}`, {
        body: data.content,
        headers: {
          "Content-Type": "application/json",
        },
      })
      .json<{ message: string }>();
  },

  renameFile: async (
    id: string,
    data: { name: string }
  ): Promise<FileMetadata> => {
    return api.put(`files/${id}/rename`, { json: data }).json<FileMetadata>();
  },

  deleteFile: async (id: string): Promise<{ message: string }> => {
    return api.delete(`files/${id}`).json<{ message: string }>();
  },
};
