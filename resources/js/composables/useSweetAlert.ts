import Swal, {
  SweetAlertResult,
  SweetAlertOptions
} from "sweetalert2";

export function useSweetAlert() {
  const colors = {
    background: "hsl(0 0% 3.9%)",
    foreground: "hsl(0 0% 98%)",
    primary: "hsl(0 0% 98%)",
    primaryText: "hsl(0 0% 9%)",
    destructive: "hsl(0 84% 60%)",
    destructiveText: "hsl(0 0% 98%)",
    border: "hsl(0 0% 14.9%)",
  };

  const baseOptions: SweetAlertOptions = {
    background: colors.background,
    color: colors.foreground,
    width: "28rem",
    padding: "1.5rem",
    position: "center",
    allowOutsideClick: true,
    allowEscapeKey: true,
    allowEnterKey: true,
  };

  const confirmAction = (
    title: string,
    text: string,
    confirmButtonText = "تأكيد"
  ): Promise<SweetAlertResult<any>> => {
    return Swal.fire({
      ...baseOptions,
      icon: "warning",
      title,
      text,
      showCancelButton: true,
      confirmButtonText,
      cancelButtonText: "الغاء",
      confirmButtonColor: colors.primary,
      cancelButtonColor: colors.border,
    });
  };

  const confirmDelete = (
    title = "هل تريد حذف العنصر؟",
    text = "لن يمكنك استعادته بعد الحذف!"
  ) => {
    return Swal.fire({
      ...baseOptions,
      icon: "warning",
      title,
      text,
      showCancelButton: true,
      confirmButtonText: "حذف",
      cancelButtonText: "الغاء",
      confirmButtonColor: colors.destructive,
      cancelButtonColor: colors.border,
    });
  };

  const confirmWarning = (title: string, text: string) => {
    return confirmAction(title, text, "متابعة");
  };

  const showSuccess = (title: string, text: string) => {
    return Swal.fire({
      ...baseOptions,
      icon: "success",
      title,
      text,
      confirmButtonText: "حسناً",
      confirmButtonColor: colors.primary,
    });
  };

  const showError = (title: string, text: string) => {
    return Swal.fire({
      ...baseOptions,
      icon: "error",
      title,
      text,
      confirmButtonText: "اغلاق",
      confirmButtonColor: colors.destructive,
    });
  };

  return {
    confirmAction,
    confirmDelete,
    confirmWarning,
    showSuccess,
    showError,
  };
}
