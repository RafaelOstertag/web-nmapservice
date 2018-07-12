pipeline {
    agent "master"

    options {
        ansiColor('xterm')
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '5')
    }

    stages {
        stage('tests') {
            steps {
                sh 'make -C backend'
            }
        }

        stage('build') {
            steps {
                sh 'make -C backend tests'
            }
        }

        stage()
    }
}