import { FileRepository } from "../../domain/repositories/FileRepository";

export class RenameItemUseCase {
  constructor(private fileRepository: FileRepository) {}

  async execute(id: string, newName: string): Promise<void> {
    return this.fileRepository.renameItem(id, newName);
  }
}
