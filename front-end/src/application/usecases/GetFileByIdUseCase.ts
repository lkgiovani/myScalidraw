import { FileEntity } from "../../domain/entities/FileEntity";
import { FileRepository } from "../../domain/repositories/FileRepository";

export class GetFileByIdUseCase {
  constructor(private fileRepository: FileRepository) {}

  async execute(id: string): Promise<FileEntity | undefined> {
    return this.fileRepository.getFileById(id);
  }
}
