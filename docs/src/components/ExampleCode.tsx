import CodeBlock from "@theme/CodeBlock";
import React from "react";

export const ExampleCode = ({
  replace,
  children,
}: {
  replace: { [key: string]: string };
  children: any;
}) => {
  for (let replacement in replace) {
    if (!children.includes(replacement)) {
      throw new Error("text does not contain replacement " + replacement);
    }
    children = children.replaceAll(replacement, replace[replacement]);
  }
  children = children
    .split("\n")
    .filter((c) => !c.trim().startsWith("_ ="))
    .map((c) => c.replaceAll("\t", "    "))
    .join("\n");
  return <CodeBlock language="go">{children}</CodeBlock>;
};
