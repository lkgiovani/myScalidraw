import { useFileStore } from "@/stores/fileStore";
import { Button } from "@/components/ui/button";
import { Icons } from "@/components/ui/icons";
import { cn } from "@/lib/utils";

export function Breadcrumb() {
  const selectedFileId = useFileStore((state) => state.selectedFileId);
  const files = useFileStore((state) => state.files);
  const selectFile = useFileStore((state) => state.selectFile);

  const buildPath = (
    fileId: string | null
  ): Array<{ id: string; name: string; isLast: boolean }> => {
    if (!fileId) return [];

    const path: Array<{ id: string; name: string; isLast: boolean }> = [];
    let currentFile = files[fileId];

    while (currentFile) {
      path.unshift({
        id: currentFile.id,
        name: currentFile.name,
        isLast: currentFile.id === fileId,
      });

      if (currentFile.parentId) {
        currentFile = files[currentFile.parentId];
      } else {
        break;
      }
    }

    return path;
  };

  const pathItems = buildPath(selectedFileId);

  return (
    <div className="flex items-center h-12 px-4 border-b border-explorer-border bg-background">
      <div className="flex items-center space-x-1 overflow-hidden">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => selectFile(null)}
          className="h-8 px-2 text-muted-foreground hover:text-foreground"
        >
          <Icons.home className="w-4 h-4" />
        </Button>

        {pathItems.length > 0 && (
          <>
            <Icons.chevronRight className="w-4 h-4 text-muted-foreground flex-shrink-0" />
            {pathItems.map((item, index) => (
              <div key={item.id} className="flex items-center space-x-1">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => selectFile(item.id)}
                  className={cn(
                    "h-8 px-2 text-sm max-w-[200px] truncate",
                    item.isLast
                      ? "text-foreground font-medium cursor-default hover:bg-transparent"
                      : "text-muted-foreground hover:text-foreground"
                  )}
                  disabled={item.isLast}
                >
                  {item.name}
                </Button>

                {index < pathItems.length - 1 && (
                  <Icons.chevronRight className="w-4 h-4 text-muted-foreground flex-shrink-0" />
                )}
              </div>
            ))}
          </>
        )}
      </div>
    </div>
  );
}
