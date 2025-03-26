import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type WtDialogProps = {
  openButtonTitle: React.ReactNode;
  form: React.ReactNode;
  onSubmitButtonClick: () => void;
  onSubmitButtonTitle: string;
  title: string;
  description?: string;
};
const WtDialog = ({ openButtonTitle, form, title, description, onSubmitButtonTitle, onSubmitButtonClick }: WtDialogProps) => {
  return (
    <Dialog>
      <DialogTrigger className={buttonVariants({ variant: "default" })}>{openButtonTitle}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>
            {description}
          </DialogDescription>
        </DialogHeader>
        {form}
        <DialogFooter>
          <DialogClose className={cn(buttonVariants({ variant: "default" }), "bg-red-500", "hover:bg-red-700")}>Cancel</DialogClose>
          <DialogClose asChild><Button onClick={() => onSubmitButtonClick()}>{onSubmitButtonTitle}</Button></DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export default WtDialog;
