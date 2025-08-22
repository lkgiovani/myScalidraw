import { useFileSystemPresenter } from "../adapters/presenters/FileSystemPresenter";
import { fileController } from "../di/container";

export const useFileSystemAdapter = () => {
  return useFileSystemPresenter(fileController);
};
