import { ExcalidrawData } from "../../domain/entities/ExcalidrawEntity";
import {
  FileEntity,
  FileSystemItemEntity,
} from "../../domain/entities/FileEntity";
import { CreateFileUseCase } from "../../application/usecases/CreateFileUseCase";
import { CreateFolderUseCase } from "../../application/usecases/CreateFolderUseCase";
import { DeleteItemUseCase } from "../../application/usecases/DeleteItemUseCase";
import { GetFileByIdUseCase } from "../../application/usecases/GetFileByIdUseCase";
import { GetFilesUseCase } from "../../application/usecases/GetFilesUseCase";
import { LoadExcalidrawFileUseCase } from "../../application/usecases/LoadExcalidrawFileUseCase";
import { RenameItemUseCase } from "../../application/usecases/RenameItemUseCase";
import { ToggleFolderUseCase } from "../../application/usecases/ToggleFolderUseCase";
import { UpdateFileUseCase } from "../../application/usecases/UpdateFileUseCase";

export class FileController {
  constructor(
    private getFilesUseCase: GetFilesUseCase,
    private getFileByIdUseCase: GetFileByIdUseCase,
    private createFileUseCase: CreateFileUseCase,
    private createFolderUseCase: CreateFolderUseCase,
    private updateFileUseCase: UpdateFileUseCase,
    private renameItemUseCase: RenameItemUseCase,
    private deleteItemUseCase: DeleteItemUseCase,
    private toggleFolderUseCase: ToggleFolderUseCase,
    private loadExcalidrawFileUseCase: LoadExcalidrawFileUseCase
  ) {}

  async getFiles(): Promise<FileSystemItemEntity[]> {
    return this.getFilesUseCase.execute();
  }

  async getFileById(id: string): Promise<FileEntity | undefined> {
    return this.getFileByIdUseCase.execute(id);
  }

  async createFile(name: string, parentId?: string): Promise<string> {
    return this.createFileUseCase.execute(name, parentId);
  }

  async createFolder(name: string, parentId?: string): Promise<string> {
    return this.createFolderUseCase.execute(name, parentId);
  }

  async updateFile(id: string, data: ExcalidrawData): Promise<void> {
    return this.updateFileUseCase.execute(id, data);
  }

  async renameItem(id: string, newName: string): Promise<void> {
    return this.renameItemUseCase.execute(id, newName);
  }

  async deleteItem(id: string): Promise<void> {
    return this.deleteItemUseCase.execute(id);
  }

  async toggleFolder(id: string): Promise<void> {
    return this.toggleFolderUseCase.execute(id);
  }

  async loadExcalidrawFile(filename: string): Promise<ExcalidrawData | null> {
    return this.loadExcalidrawFileUseCase.execute(filename);
  }
}
