import { ExcalidrawData } from "../../domain/entities/ExcalidrawEntity";
import { FileRepository } from "../../domain/repositories/FileRepository";

export class LoadExcalidrawFileUseCase {
  constructor(private fileRepository: FileRepository) {}

  async execute(filename: string): Promise<ExcalidrawData | null> {
    return this.fileRepository.loadExcalidrawFile(filename);
  }
}
