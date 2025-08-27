import { useFileStore } from "@/stores/fileStore";
import { ChevronDown, ChevronRight, FileText } from "lucide-react";
import { Button } from "@/components/ui/button";
import { FileContextMenu } from "@/components/molecules/FileContextMenu";
import { cn } from "@/lib/utils";

interface FileTreeNodeProps {
  fileId: string;
  depth?: number;
}

function FileTreeNode({ fileId, depth = 0 }: FileTreeNodeProps) {
  const file = useFileStore((state) => state.files[fileId]);
  const isExpanded = useFileStore((state) => state.expandedFolders.has(fileId));
  const selectedFileId = useFileStore((state) => state.selectedFileId);
  const toggleFolder = useFileStore((state) => state.toggleFolder);
  const selectFile = useFileStore((state) => state.selectFile);

  if (!file) return null;

  const hasChildren =
    file.type === "folder" && file.children && file.children.length > 0;

  const isSelected = selectedFileId === fileId;

  const handleClick = () => {
    if (file.type === "folder") {
      toggleFolder(fileId);
    } else {
      selectFile(fileId);
    }
  };

  return (
    <div className="w-full">
      <FileContextMenu file={file}>
        <Button
          variant="ghost"
          onClick={handleClick}
          className={cn(
            "w-full justify-start h-8 px-2 py-1 text-sm font-normal",
            "hover:bg-explorer-sidebar-hover transition-colors",
            isSelected && "bg-explorer-sidebar-active text-primary",
            "focus:bg-explorer-sidebar-hover focus:outline-none"
          )}
          style={{ paddingLeft: `${depth * 16 + 8}px` }}
        >
          <div className="flex items-center gap-2 min-w-0 flex-1">
            {file.type === "folder" && (
              <div className="flex-shrink-0 w-4 h-4">
                {hasChildren ? (
                  isExpanded ? (
                    <ChevronDown className="w-3 h-3 text-muted-foreground" />
                  ) : (
                    <ChevronRight className="w-3 h-3 text-muted-foreground" />
                  )
                ) : null}
              </div>
            )}

            <FileText
              className={cn(
                "w-4 h-4 flex-shrink-0",
                file.type === "folder"
                  ? "text-file-folder"
                  : "text-file-document"
              )}
            />

            <span className="truncate text-left">
              {file.name || `[Nome vazio - ID: ${fileId.slice(0, 8)}...]`}
            </span>
          </div>
        </Button>
      </FileContextMenu>

      {file.type === "folder" && isExpanded && hasChildren && (
        <div className="ml-0">
          {file.children?.map((childId) => (
            <FileTreeNode key={childId} fileId={childId} depth={depth + 1} />
          ))}
        </div>
      )}
    </div>
  );
}

export function FileTree() {
  const rootFolders = useFileStore((state) => state.rootFolders);

  return (
    <div className="w-full space-y-1 p-2">
      {rootFolders.map((folderId) => (
        <FileTreeNode key={folderId} fileId={folderId} />
      ))}
    </div>
  );
}
