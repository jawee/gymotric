import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type WtDialogProps = {
  openButtonTitle?: React.ReactNode;
  openButtonClassName?: string;
  form: React.ReactNode;
  onSubmitButtonClick: () => void;
  onSubmitButtonTitle: React.ReactNode;
  title: string;
  description?: string;
  dialogProps?: React.ComponentProps<typeof Dialog>;
  shouldUseDefaultSubmit?: boolean;
  topPercentage?: string;
  onOpenAutoFocus?: (e: Event) => void;
};
const WtDialog = ({ openButtonTitle, openButtonClassName, form, title, description, onSubmitButtonTitle, onSubmitButtonClick, shouldUseDefaultSubmit = true, topPercentage = "20", onOpenAutoFocus, dialogProps }: WtDialogProps) => {
  let topPosition = "max-md:top-[20%]";
  if (topPercentage == "25") {
    topPosition = "max-md:top-[25%]";
  }

  if (openButtonClassName === null || openButtonClassName === undefined) {
    openButtonClassName = buttonVariants({ variant: "default" });
  }

  return (
    <Dialog {...dialogProps} >
      {openButtonTitle !== null && openButtonTitle !== undefined ? <DialogTrigger className={openButtonClassName}>{openButtonTitle}</DialogTrigger> : <></>}
      <DialogContent className={topPosition} onOpenAutoFocus={onOpenAutoFocus}>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>
            {description}
          </DialogDescription>
        </DialogHeader>
        {form}
        <DialogFooter>
          <DialogClose className={cn(buttonVariants({ variant: "default" }), "bg-red-500", "hover:bg-red-700")}>Cancel</DialogClose>
          {shouldUseDefaultSubmit ?
            <DialogClose asChild><Button onClick={() => onSubmitButtonClick()}>{onSubmitButtonTitle}</Button></DialogClose> :
            <Button onClick={onSubmitButtonClick}>Submit</Button>
          }
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export default WtDialog;
