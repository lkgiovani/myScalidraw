import {
  Folder,
  FolderOpen,
  File,
  FileText,
  Image,
  Code,
  Search,
  ChevronRight,
  ChevronDown,
  Home,
  MoreHorizontal,
} from "lucide-react";

export const Icons = {
  folder: Folder,
  folderOpen: FolderOpen,
  file: File,
  fileText: FileText,
  image: Image,
  code: Code,
  search: Search,
  chevronRight: ChevronRight,
  chevronDown: ChevronDown,
  home: Home,
  moreHorizontal: MoreHorizontal,
};

interface FileIconProps {
  type: "file" | "folder";
  name?: string;
  isOpen?: boolean;
  className?: string;
}

export function FileIcon({ type, name, isOpen, className }: FileIconProps) {
  if (type === "folder") {
    const Icon = isOpen ? Icons.folderOpen : Icons.folder;
    return <Icon className={className} />;
  }

  const extension = name?.split(".").pop()?.toLowerCase();

  switch (extension) {
    case "excalidraw":
      return <Icons.fileText className={className} />;
    case "png":
    case "jpg":
    case "jpeg":
    case "gif":
    case "svg":
      return <Icons.image className={className} />;
    case "js":
    case "ts":
    case "tsx":
    case "jsx":
    case "html":
    case "css":
      return <Icons.code className={className} />;
    default:
      return <Icons.file className={className} />;
  }
}
