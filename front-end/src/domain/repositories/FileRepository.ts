import {
  FileEntity,
  FileSystemItemEntity,
  FolderEntity,
} from "../entities/FileEntity";
import { ExcalidrawData } from "../entities/ExcalidrawEntity";

export interface FileRepository {
  getFiles(): Promise<FileSystemItemEntity[]>;
  getFileById(id: string): Promise<FileEntity | undefined>;
  createFile(name: string, parentId?: string): Promise<string>;
  createFolder(name: string, parentId?: string): Promise<string>;
  updateFile(id: string, data: ExcalidrawData): Promise<void>;
  renameItem(id: string, newName: string): Promise<void>;
  deleteItem(id: string): Promise<void>;
  toggleFolder(id: string): Promise<void>;
  loadExcalidrawFile(filename: string): Promise<ExcalidrawData | null>;
}
