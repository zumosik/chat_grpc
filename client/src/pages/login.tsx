import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { useState } from "react"

import { Button } from "@/components/ui/button"
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { useToast } from "@/components/ui/use-toast"
import { useNavigate } from "react-router-dom"

const formSchema = z.object({
    username: z.string().min(5, {
        message: "Username must be at least 5 characters.",
    }),
    email: z.string().email(
        "Please enter a valid email address."
    ),
    password: z.string().min(8, {
        message: "Password must be at least 8 characters."
    })
})



export default function LoginPage({ isLogin }: { isLogin: boolean }) {
    const { toast } = useToast()
    const navigate = useNavigate();


    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            username: "",
            email: "",
            password: "",
        },
    })


    function onSubmit(values: z.infer<typeof formSchema>) {
        console.log(values) // TODO: send to server

        if (isLoginProp) {
            toast({
                title: "Logged in",
                description: "You have been logged in successfully."
            })
        } else {
            toast({
                title: "Account created",
                description: "Your account has been created successfully, now you need to confirm your email address."
            })
        }

        navigate("/");
    }

    const [isLoginProp, setIsLoginProp] = useState(isLogin)


    return (
        <div className="flex justify-center items-center h-screen">
            <div>
                <h1 className="text-3xl mb-10">
                    {isLoginProp ? "Login into system" : "Create new account"}
                </h1>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                        <FormField
                            control={form.control}
                            name="username"
                            render={({ field }) => (
                                <FormItem>
                                    {/* <FormLabel>Username</FormLabel> */}
                                    <FormControl>
                                        <Input placeholder="Username" {...field} />
                                    </FormControl>
                                    <FormDescription>
                                        Name that other people will see.
                                    </FormDescription>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="email"
                            render={({ field }) => (
                                <FormItem>
                                    {/* <FormLabel>Email</FormLabel> */}
                                    <FormControl>
                                        <Input placeholder="Email" {...field} />
                                    </FormControl>
                                    <FormDescription>
                                        Your real email address.
                                    </FormDescription>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="password"
                            render={({ field }) => (
                                <FormItem>
                                    {/* <FormLabel>Password</FormLabel> */}
                                    <FormControl>
                                        <Input type="password"  placeholder="Password" {...field} />
                                    </FormControl>
                                    <FormDescription>
                                        Your strong and secure password.
                                    </FormDescription>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <Button type="submit" className="w-full">Submit</Button>
                    </form>


                    <span className="inline-block mt-4 cursor-pointer">
                        {isLoginProp ? "Don't have an account?  " : "Already have an account?  "}
                        <button>
                            <span onClick={() => {
                                setIsLoginProp(!isLoginProp)
                                form.reset()
                            }} className="text-blue-500"> {isLoginProp ? "Register" : "Login"}</span>
                        </button>

                    </span>


                </Form>

            </div>
        </div>
    )
}
