import Editor from "@monaco-editor/react";


const defaultOptions = {
    selectOnLineNumbers: true,
    roundedSelection: false,
    readOnly: false,
    cursorStyle: 'line',
    theme: 'vs-dark',
    minimap: {
        enabled: false
    }
};

const readOnlyOptions = {
    selectOnLineNumbers: true,
    roundedSelection: false,
    readOnly: true,
    cursorStyle: 'line',
    theme: 'vs-dark',
    minimap: {
        enabled: false
    }
}

export const Console = ({ children, readOnly, code = '', setCode }) => {
    let options = readOnly ? readOnlyOptions : defaultOptions;
    let title = readOnly ? "Consola" : "Entrada";
    return (
        <div className='col d-flex flex-column justify-content-evenly'>
            
            {children}
            <h4 className="text-light">{title}</h4>

            <Editor
                theme="vs-dark"
                // defaultLanguage="rust"
                value={code}
                onChange={setCode}
                options = {options}
            />
        </div>
    )
}


