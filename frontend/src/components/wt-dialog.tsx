import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type WtDialogProps = {
  openButtonTitle?: React.ReactNode;
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
const WtDialog = ({ openButtonTitle, form, title, description, onSubmitButtonTitle, onSubmitButtonClick, shouldUseDefaultSubmit = true, topPercentage = "20", onOpenAutoFocus, dialogProps }: WtDialogProps) => {
  let topPosition = "max-md:top-[20%]";
  if (topPercentage == "25") {
    topPosition = "max-md:top-[25%]";
  }

  return (
    <Dialog {...dialogProps} >
      {openButtonTitle === null ? <DialogTrigger className={buttonVariants({ variant: "default" })}>{openButtonTitle}</DialogTrigger> : <></>}
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
