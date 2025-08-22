import { FileRepository } from "../../domain/repositories/FileRepository";

export class ToggleFolderUseCase {
  constructor(private fileRepository: FileRepository) {}

  async execute(id: string): Promise<void> {
    return this.fileRepository.toggleFolder(id);
  }
}
