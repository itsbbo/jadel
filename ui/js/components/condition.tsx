import React, { ReactElement, ReactNode } from 'react';

interface Props {
    condition: boolean;
    children: ReactNode;
}

const If: React.FC<Props> = ({ condition, children }) => {
    let trueBranch: ReactElement | null = null;
    let falseBranch: ReactElement | null = null;

    React.Children.forEach(children, (child) => {
        if (React.isValidElement(child)) {
            const childType = (child.type as any).displayName;
            if (childType === 'Else') {
                falseBranch = child;
            } else {
                trueBranch = child;
            }
        }
    });

    return condition ? trueBranch : falseBranch;
};

const Else: React.FC<{ children: ReactNode }> = ({ children }) => <>{children}</>;

If.displayName = 'If';
Else.displayName = 'Else';

export { Else, If };
