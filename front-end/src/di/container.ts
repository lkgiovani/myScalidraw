import {
  CreateFileUseCase,
  CreateFolderUseCase,
  DeleteItemUseCase,
  GetFileByIdUseCase,
  GetFilesUseCase,
  LoadExcalidrawFileUseCase,
  RenameItemUseCase,
  ToggleFolderUseCase,
  UpdateFileUseCase,
} from "../application/usecases";
import { FileRepository } from "../domain/repositories/FileRepository";
import { ApiFileRepository } from "../infrastructure/repositories/ApiFileRepository";
import { FileController } from "../adapters/controllers/FileController";

// Escolha qual repositório usar
const useApi = true; // Mude para false para usar o repositório local

// Criar instância do repositório
let fileRepository: FileRepository;

if (useApi) {
  fileRepository = new ApiFileRepository("http://localhost:8081/api");
}
// Criar use cases
const getFilesUseCase = new GetFilesUseCase(fileRepository);
const getFileByIdUseCase = new GetFileByIdUseCase(fileRepository);
const createFileUseCase = new CreateFileUseCase(fileRepository);
const createFolderUseCase = new CreateFolderUseCase(fileRepository);
const updateFileUseCase = new UpdateFileUseCase(fileRepository);
const renameItemUseCase = new RenameItemUseCase(fileRepository);
const deleteItemUseCase = new DeleteItemUseCase(fileRepository);
const toggleFolderUseCase = new ToggleFolderUseCase(fileRepository);
const loadExcalidrawFileUseCase = new LoadExcalidrawFileUseCase(fileRepository);

// Criar controller
const fileController = new FileController(
  getFilesUseCase,
  getFileByIdUseCase,
  createFileUseCase,
  createFolderUseCase,
  updateFileUseCase,
  renameItemUseCase,
  deleteItemUseCase,
  toggleFolderUseCase,
  loadExcalidrawFileUseCase
);

export { fileController };
