import { ExcalidrawData } from "../../domain/entities/ExcalidrawEntity";
import {
  FileEntity,
  FileSystemItemEntity,
  FolderEntity,
} from "../../domain/entities/FileEntity";
import { FileRepository } from "../../domain/repositories/FileRepository";

interface ApiFileResponse {
  id: string;
  name: string;
  isFolder?: boolean;
  children?: ApiFileResponse[];
  data?: Record<string, unknown>;
  lastModified?: number;
  parentId?: string;
  isExpanded?: boolean;
}

export class ApiFileRepository implements FileRepository {
  private apiBaseUrl: string;
  private cachedFiles: FileSystemItemEntity[] = [];

  constructor(apiBaseUrl: string = "http://localhost:8081/api") {
    this.apiBaseUrl = apiBaseUrl;
  }

  async getFiles(): Promise<FileSystemItemEntity[]> {
    try {
      const response = await fetch(`${this.apiBaseUrl}/files`);
      if (!response.ok) {
        throw new Error(`Failed to fetch files: ${response.statusText}`);
      }

      const data = (await response.json()) as ApiFileResponse[];
      this.cachedFiles = this.transformApiFiles(data);
      return this.cachedFiles;
    } catch (error) {
      console.error("Error fetching files:", error);
      return [];
    }
  }

  async getFileById(id: string): Promise<FileEntity | undefined> {
    try {
      const response = await fetch(`${this.apiBaseUrl}/files/${id}`);
      if (!response.ok) {
        throw new Error(`Failed to fetch file ${id}: ${response.statusText}`);
      }

      const data = (await response.json()) as ApiFileResponse;
      return this.transformApiFile(data) as FileEntity;
    } catch (error) {
      console.error(`Error fetching file ${id}:`, error);
      return undefined;
    }
  }

  async createFile(name: string, parentId?: string): Promise<string> {
    // Simulando criação de arquivo - em uma implementação real, isso seria uma chamada POST
    const newId = `file-${Date.now()}`;
    console.log(
      `Creating file ${name} with ID ${newId} under parent ${
        parentId || "root"
      }`
    );

    // Em uma implementação real, atualizaríamos o cache com a resposta da API
    return newId;
  }

  async createFolder(name: string, parentId?: string): Promise<string> {
    // Simulando criação de pasta - em uma implementação real, isso seria uma chamada POST
    const newId = `folder-${Date.now()}`;
    console.log(
      `Creating folder ${name} with ID ${newId} under parent ${
        parentId || "root"
      }`
    );

    // Em uma implementação real, atualizaríamos o cache com a resposta da API
    return newId;
  }

  async updateFile(id: string, data: ExcalidrawData): Promise<void> {
    try {
      // Em uma implementação real, isso seria uma chamada PUT ou PATCH
      console.log(`Updating file ${id} with new data`);

      // Simulando uma chamada à API para salvar o desenho
      await fetch(`${this.apiBaseUrl}/sala/${id}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });
    } catch (error) {
      console.error(`Error updating file ${id}:`, error);
    }
  }

  async renameItem(id: string, newName: string): Promise<void> {
    // Simulando renomeação - em uma implementação real, isso seria uma chamada PATCH
    console.log(`Renaming item ${id} to ${newName}`);
  }

  async deleteItem(id: string): Promise<void> {
    // Simulando exclusão - em uma implementação real, isso seria uma chamada DELETE
    console.log(`Deleting item ${id}`);
  }

  async toggleFolder(id: string): Promise<void> {
    // Esta operação é apenas do lado do cliente, não precisa de API
    console.log(`Toggling folder ${id}`);
  }

  async loadExcalidrawFile(filename: string): Promise<ExcalidrawData | null> {
    try {
      const response = await fetch(`${this.apiBaseUrl}/excalidraw`);
      if (!response.ok) {
        throw new Error(`Failed to load ${filename}`);
      }
      return (await response.json()) as ExcalidrawData;
    } catch (error) {
      console.error("Error loading Excalidraw file:", error);
      return null;
    }
  }

  // Métodos auxiliares para transformar dados da API no formato do domínio
  private transformApiFiles(
    apiFiles: ApiFileResponse[]
  ): FileSystemItemEntity[] {
    return apiFiles.map((file) => this.transformApiFile(file));
  }

  private transformApiFile(apiFile: ApiFileResponse): FileSystemItemEntity {
    if (apiFile.isFolder) {
      return {
        id: apiFile.id,
        name: apiFile.name,
        isFolder: true,
        parentId: apiFile.parentId,
        children: apiFile.children
          ? this.transformApiFiles(apiFile.children)
          : [],
        isExpanded: apiFile.isExpanded || false,
      } as FolderEntity;
    } else {
      return {
        id: apiFile.id,
        name: apiFile.name,
        data: (apiFile.data as ExcalidrawData) || {
          elements: [],
          appState: {
            viewBackgroundColor: "#ffffff",
            currentItemStrokeColor: "#000000",
            currentItemBackgroundColor: "#ffffff",
            currentItemFillStyle: "solid",
            currentItemStrokeWidth: 1,
            currentItemStrokeStyle: "solid",
            currentItemRoughness: 1,
            currentItemOpacity: 100,
            currentItemFontFamily: 1,
            currentItemFontSize: 20,
            currentItemTextAlign: "left",
            currentItemStartArrowhead: null,
            currentItemEndArrowhead: null,
            scrollX: 0,
            scrollY: 0,
            zoom: { value: 1 },
            currentItemRoundness: "round",
            gridSize: null,
            colorPalette: {},
          },
        },
        lastModified: apiFile.lastModified || Date.now(),
        parentId: apiFile.parentId,
      } as FileEntity;
    }
  }
}
