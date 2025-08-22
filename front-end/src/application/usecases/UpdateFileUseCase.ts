import { ExcalidrawData } from "../../domain/entities/ExcalidrawEntity";
import { FileRepository } from "../../domain/repositories/FileRepository";

export class UpdateFileUseCase {
  constructor(private fileRepository: FileRepository) {}

  async execute(id: string, data: ExcalidrawData): Promise<void> {
    return this.fileRepository.updateFile(id, data);
  }
}
