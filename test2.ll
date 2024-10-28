@.str = private unnamed_addr constant [4 x i8] c"%d\0A\00", align 1
declare i32 @printf(i8*, ...)
define i32 @main() {
entry:
%0 = add i32 1, 2
%1 = add i32 %0, 3
%2 = getelementptr [4 x i8], [4 x i8]* @.str, i32 0, i32 0
call i32 (i8*, ...) @printf(i8* %2, i32 %1)
ret i32 0
}
