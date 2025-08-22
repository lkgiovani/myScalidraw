import { FileRepository } from "../../domain/repositories/FileRepository";

export class CreateFileUseCase {
  constructor(private fileRepository: FileRepository) {}

  async execute(name: string, parentId?: string): Promise<string> {
    return this.fileRepository.createFile(name, parentId);
  }
}
