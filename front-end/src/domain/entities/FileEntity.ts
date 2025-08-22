import { ExcalidrawData } from "./ExcalidrawEntity";

export interface FileEntity {
  id: string;
  name: string;
  data: ExcalidrawData;
  lastModified: number;
  parentId?: string;
}

export interface FolderEntity {
  id: string;
  name: string;
  isFolder: true;
  parentId?: string;
  children: (FileEntity | FolderEntity)[];
  isExpanded?: boolean;
}

export type FileSystemItemEntity = FileEntity | FolderEntity;
