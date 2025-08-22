import { FileRepository } from "../../domain/repositories/FileRepository";

export class CreateFolderUseCase {
  constructor(private fileRepository: FileRepository) {}

  async execute(name: string, parentId?: string): Promise<string> {
    return this.fileRepository.createFolder(name, parentId);
  }
}
