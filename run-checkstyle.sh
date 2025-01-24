#!/bin/sh

# Define the URL for the Checkstyle JAR file
CHECKSTYLE_VERSION="8.42"
CHECKSTYLE_JAR="checkstyle-${CHECKSTYLE_VERSION}-all.jar"
CHECKSTYLE_URL="https://github.com/checkstyle/checkstyle/releases/download/checkstyle-${CHECKSTYLE_VERSION}/${CHECKSTYLE_JAR}"

# Define the directory where the JAR file should be stored
LIB_DIR="libs"

# Create the directory if it does not exist
mkdir -p "$LIB_DIR"

# Define the full path to the JAR file
CHECKSTYLE_PATH="${LIB_DIR}/${CHECKSTYLE_JAR}"

# Check if the JAR file already exists
if [ -f "$CHECKSTYLE_PATH" ]; then
    echo "Checkstyle JAR file already exists at ${CHECKSTYLE_PATH}"
else
    echo "Checkstyle JAR file not found. Downloading from ${CHECKSTYLE_URL}..."
    curl -L -o "$CHECKSTYLE_PATH" "$CHECKSTYLE_URL" && echo "downloaded successfully"
fi

# Add the lib directory to .gitignore if it's not already present
GITIGNORE_FILE=".gitignore"

if ! grep -q "^${LIB_DIR}/$" "$GITIGNORE_FILE"; then
    echo "Adding ${LIB_DIR}/ to ${GITIGNORE_FILE}"
    echo "${LIB_DIR}/" >> "$GITIGNORE_FILE"
else
    echo "${LIB_DIR}/ is already in ${GITIGNORE_FILE}"
fi

# Check if there are any Java files in the project
JAVA_FILES_FOUND=$(find . -name "*.java")

# if there is java file go ahead to trigger checkStyle for all java files
if [ -z "$JAVA_FILES_FOUND" ]; then
    echo "No Java files found in the project."
    exit 0
else
    echo "Java files found in the project. Now run further checks."
fi

# Path to your Checkstyle configuration file
CHECKSTYLE_CONFIG="./conf/checkstyle/checkStyleAll.xml"

# shellcheck disable=SC2027
# shellcheck disable=SC2046

All_Java_Files=$(find . -name "*.java")

# Run Checkstyle on all Java files in the repository
java -jar "$CHECKSTYLE_PATH" -c "$CHECKSTYLE_CONFIG" "$All_Java_Files"

echo "finish ...."

# Capture the exit code of Checkstyle
# shellcheck disable=SC2320
STATUS=$?

# Exit with the same status code as Checkstyle
exit $STATUS
