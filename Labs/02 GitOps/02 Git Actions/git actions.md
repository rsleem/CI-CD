## <font color='red'> 2.2 GitHub Actions </font>
GitHub Actions help you automate tasks within your software development life cycle. GitHub Actions are event-driven, meaning that you can run a series of commands after a specified event has occurred. For example, every time someone creates a pull request for a repository, you can automatically run a command that executes a software testing script.

  > For further information: https://docs.github.com/en/actions

In this lab we're going to:
* create a simple GitHub Action

  > MarketPlace for Actions: https://github.com/marketplace?type=actions


---


**Pre-requisties:**
* GitHub Account
If you want to try these GitHub actions you will require a GitHub account.

  > https://github.com/join

Don't worry if you dont have access as the videos will guide you through the Workflows.

---

#### <font color='red'>2.2.1 GitHub Actions</font>
You can set up continuous integration for your project using a workflow template that matches the language and tooling you want to use.

List of GitHub Actions:
* github actions demo
* hello world
* different shells


So lets dive in the deep end and create a workflow thats shows whats happening behind the scenes.

  > GitHub Actions repository:  http://github.com/jporeilly/GitHub-Actions.git

**github actions demo**
* create a new branch for the workflow: github-actions-demo
* then add this workflow - github-actions-demo.yaml - to the github-actions-demo/.github/workflows/  directory:

```
name: GitHub Actions Demo
on: [push]
jobs:
  Explore-GitHub-Actions:
    runs-on: ubuntu-latest
    steps:
      - run: echo "🎉 The job was automatically triggered by a ${{ github.event_name }} event."
      - run: echo "🐧 This job is now running on a ${{ runner.os }} server hosted by GitHub!"
      - run: echo "🔎 The name of your branch is ${{ github.ref }} and your repository is ${{ github.repository }}."
      - name: Check out repository code
        uses: actions/checkout@v2
      - run: echo "💡 The ${{ github.repository }} repository has been cloned to the runner."
      - run: echo "🖥️ The workflow is now ready to test your code on the runner."
      - name: List files in the repository
        run: |
          ls ${{ github.workspace }}
      - run: echo "🍏 This job's status is ${{ job.status }}."
```

* This will trigger a push event.
* On GitHub, navigate to the main page of the repository.
* Under your repository name, click Actions.
* In the left sidebar, click the workflow you want to see.
* From the list of workflow runs, click the name of the run you want to see.
* Under Jobs, click the Explore-GitHub-Actions job.

**hello world**
* create a new branch for the workflow: hello-world
* then add this workflow - hello-world.yaml - to the hello-world/.github/workflows/  directory:

```
name: 01 Hello World 
on: [push]
jobs:
  run-shell-command:
    runs-on: ubuntu-latest
    steps: 
      - name: Echo a string
        run: echo "Hello World"
      - name: Multiline script 
        run: |
           node -v 
           npm -v
```
* This will trigger a push event.
* On GitHub, navigate to the main page of the repository.
* Under your repository name, click Actions.


**different shells**
* create a new branch for the workflow: different-shells
* then add this workflow - different-shells.yaml - to the different-shells/.github/workflows/  directory:

```
name: 02 Different Shells 
on: [push]
jobs:
  run-shell-command:
    runs-on: ubuntu-latest
    steps: 
      - name: echo a string
        run: echo "Hello World"
      - name: multiline script 
        run: |
           node -v 
           npm -v
      - name: python Command 
        run: |
          import platform 
          print
          (platform.processor())
        shell: python
  run-windows-command:
    runs-on: windows-latest
    needs: ["run-shell-command"]
    steps:
      - name: Directory PowerShell
        run: Get-Location 
      - name: Directory Bash 
        run: pwd 
        shell: bash 
```
* This will trigger a push event.
* On GitHub, navigate to the main page of the repository.
* Under your repository name, click Actions.


**hello world - Marketplace**
* create a new branch for the workflow: marketplace
* click on the Actions tab
* click on the New Workflow
* select 'Simple Workflow'
* deploy to main

---