import { FileRepository } from "../../domain/repositories/FileRepository";

export class DeleteItemUseCase {
  constructor(private fileRepository: FileRepository) {}

  async execute(id: string): Promise<void> {
    return this.fileRepository.deleteItem(id);
  }
}
