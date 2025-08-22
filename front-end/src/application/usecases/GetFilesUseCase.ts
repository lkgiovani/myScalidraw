import { FileSystemItemEntity } from "../../domain/entities/FileEntity";
import { FileRepository } from "../../domain/repositories/FileRepository";

export class GetFilesUseCase {
  constructor(private fileRepository: FileRepository) {}

  async execute(): Promise<FileSystemItemEntity[]> {
    return this.fileRepository.getFiles();
  }
}
