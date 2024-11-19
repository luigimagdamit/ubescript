@.str = private unnamed_addr constant [4 x i8] c"%d\0A\00", align 1
declare i32 @printf(i8*, ...)
define i32 @main() {
entry:
%0 = add i32 2000, 420
%1 = add i32 %0, 1
%2 = add i32 %1, 2
%3 = sub i32 %2, 69
%4 = getelementptr [4 x i8], [4 x i8]* @.str, i32 0, i32 0
call i32 (i8*, ...) @printf(i8* %4, i32 %3)
ret i32 0
}
